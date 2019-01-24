# Status Indicator

The `ctop` grid view provides a compact status indicator to convey container state

<img width="200px" src="img/status.png" alt="ctop"/>

<span align="center">

Appearance | Description
--- | ---
red | container is stopped
green | container is running
▮▮ | container is paused

</span>

If the container is configured with a health check, a `+` will appear next to the indicator

<span align="center">

Appearance | Description
--- | ---
red | health check in failed state
yellow | health check in starting state
green | health check in OK state

</span>
