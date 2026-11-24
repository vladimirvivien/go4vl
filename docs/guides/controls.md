# Controls & Events

Hardware cameras, encoders, and capture cards expose adjustable parameters (brightness, exposure, gain, bitrate, codec profiles) alongside asynchronous driver notifications (control updates, source resolution changes, end-of-stream). `go4vl` separates single-value user controls from atomic multi-value extended controls and provides typed wrappers for V4L2 event subscriptions.

## Querying and setting user controls

User controls address scalar settings such as `CtrlBrightness`, `CtrlContrast`, `CtrlSaturation`, `CtrlHue`, and `CtrlExposure`. Call `QueryAllControls()` on an open `Device` to inspect supported control ranges, step sizes, defaults, and menu options, then read or modify values by ID:

```go
ctrls, err := dev.QueryAllControls()
if err != nil {
	log.Fatal(err)
}
for _, c := range ctrls {
	fmt.Printf("%s (id=0x%x): min=%d max=%d step=%d default=%d\n",
		c.Name, c.ID, c.Minimum, c.Maximum, c.Step, c.Default)
}

curr, err := dev.GetControl(v4l2.CtrlBrightness)
if err == nil {
	fmt.Printf("Current brightness: %d\n", curr.Value)
}
if err := dev.SetControlValue(v4l2.CtrlBrightness, 128); err != nil {
	log.Fatal(err)
}
```

Convenience helpers (`SetControlBrightness`, `SetControlContrast`, `SetControlSaturation`, and `SetControlHue`) wrap common image adjustments directly on `*device.Device`. For menu-typed controls (`c.IsMenu()`), call `c.GetMenuItems()` to list valid discrete entries. See `examples/user_ctrl/` and `examples/control_reference/`.

## Atomic extended controls (`ExtControls`)

Standard user-control ioctls update one 32-bit integer at a time. When you need to update multiple controls atomically, configure 64-bit values or strings, write camera/codec control classes together, or submit compound kernel structs, use `v4l2.ExtControls`:

```go
ctrls := v4l2.NewExtControlsWithClass(v4l2.CtrlClassCamera)
ctrls.AddValue(v4l2.CtrlExposureAuto, 1)
ctrls.AddValue(v4l2.CtrlExposureAbsolute, 250)

// Validate with the driver before applying
if err := dev.TryExtControls(ctrls); err != nil {
	log.Fatalf("try controls: %v (errIdx=%d)", err, ctrls.GetErrorIndex())
}
if err := dev.SetExtControls(ctrls); err != nil {
	log.Fatalf("set controls: %v", err)
}
```

`TryExtControls` invokes `VIDIOC_TRY_EXT_CTRLS` so the driver validates and clamps fields without mutating active hardware state; `SetExtControls` applies the full batch atomically (`VIDIOC_S_EXT_CTRLS`). When passing string or compound `[]byte` payloads, `v4l2` pins Go memory via `runtime.Pinner` across the syscall so garbage collection never moves buffers during kernel execution.

For stateless hardware decoders (H.264, HEVC, VP8, VP9, AV1, MPEG2, FWHT), call `ctrls.SetRequestFD(req.Fd())` to bind an extended control batch directly to a `MediaRequest`—see [Codecs & Memory-to-Memory](codecs-m2m.md). See `examples/ext_ctrls/`.

## Subscribing to and dequeuing V4L2 events

Drivers signal out-of-band state changes—including mid-stream resolution changes (`EventSourceChange`), control modifications (`EventCtrl`), frame sync (`EventFrameSync`), and codec drain completion (`EventEOS`)—as V4L2 events delivered over `POLLPRI`.

Subscribe to the events you care about, wait for `v4l2.PollPri` via `v4l2.Poll` or `v4l2.Poller`, and dequeue typed event payloads:

```go
sub := v4l2.NewEventSubscription(v4l2.EventSourceChange)
if err := dev.SubscribeEvent(sub); err != nil {
	log.Fatal(err)
}
defer dev.UnsubscribeEvent(v4l2.NewEventSubscription(v4l2.EventSourceChange))

evFlags, err := v4l2.Poll(ctx, dev.Fd(), v4l2.PollPri, 2*time.Second)
if err != nil {
	log.Fatal(err)
}
if evFlags&v4l2.PollPri != 0 {
	ev, err := dev.DequeueEvent()
	if err == nil && ev.GetType() == v4l2.EventSourceChange {
		sc := ev.GetSrcChangeData()
		fmt.Printf("Source changed (changes=0x%x)\n", sc.Changes)
	}
}
```

Use `v4l2.NewControlEventSubscription(ctrlID)` with `v4l2.EventSubFlagSendInitial` to receive an immediate initial value notification followed by asynchronous updates whenever another process modifies that control.
