# sig0namectl Javascript UI documentation

## Javascript Event System

The UI can listen to events in order to execute tasks
related to these events.
The following events exist.

### WASM Events

- `wasm_ready`: This event occurs, when the WASM is fully loaded and ready to use.  All functions using WASM should only be started after this event occurred.

### KEY Store Events

The Javascript key store representation signals changes via the following events:

- `keys_ready`
- `keys_updated`
