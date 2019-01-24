# Status Indicator

The `ctop` grid view provides a compact status indicator to convey container state
<p align="center"><img width="200px" src="img/status.png" alt="ctop"/></p>

Appearance | Description
--- | ---
red | container is stopped
green | container is running
two vertical bars (▮▮) | container is paused

If the container is configured with a health check, a `+` will appear next to the indicator

Appearance | Description
--- | ---
red | health check in failed state
yellow | health check in starting state
green | health check in OK state
