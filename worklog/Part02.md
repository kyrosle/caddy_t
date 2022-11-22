# Main

Continue to build the fie

# Building HTTP

## `Context` (context.go)

`Context` is a type which defines the lifetime of modules that
are loaded and provides access to the parent configuration
that spawned the modules which are loaded. It should be used
standard context package only if you don't need the Caddy
specified features. These contexts are canceled when the 
lifetime of the modules loaded from it is over.
```go
type Context struct {
	context.Context
	moduleInstances map[string][]Module
	cfg 			*Config
	cleanupFuncs 	[]func()
	ancestry 		[]Module
}
```

Use function `NewContext()` to get a valid value 
(but most modules will not actually need to do this).

__Context Fields Using Types__ :
* `Module` type :

(modules.go)

`Module` is a type that is used as a Caddy module. 
In addition to this interface, most modules will implement some 
interface expected by their host module in order to be useful.
To learn which interface(s) to implement,
see the documentation for the host module. At a bare minimum,
this interface, when implemented, only provides the module's ID and 
constructor function. 
```go
type Module interface {
	// This method indicates that the type is a Caddy module.
	// The returned ModuleInfo must have both a name and a constructor function.
	// This method must not have any side-effects.
	CaddyModule() ModuleInfo
}
```
`Module` will often implement additional interfaces
including `Provisioner`, `Validator`, and `CleanerUpper`.
If a module implements these interfaces, their methods are called
during the module's lifespan.

When a module is loaded by a host module, the following happens: 
1. `ModuleInfo.New()` is called to get a new instance of the module.
2. The module's configuration is unmarshalled into that instance.
3. If the module is a `Provisioner` the `Provision()` method is called.
4. If the module is a `Validator` the `Validate()` method is called.
5. The module will probably be type-asserted from `any` to some other, 
more useful interface expected by the host module. For example, HTTP handler
modules are type-asserted as `caddyhttp.MiddlewareHandler` values.
6. When a module's containing Context is canceled, if it is a `CleanerUpper`, 
its `Cleaner()` method is called.

`ModuleInfo` 

Represents a registered Caddy module.
```go
type ModuleInfo struct {
	ID  ModuleID
	New func() Module
}
```
__ModuleInfo Fields__ :
* `ID` : ID is the "full name" of the module. 
It must be unique and properly namespaced.
* `New` : 
New returns a pointer to a new, empty 
interface of the module's type. This
method must not have any side-effects,
and no other initialization should 
occur within it. Any initialization
of the returned value should be done
in a `Provision()` method.

`ModuleID`

Is a string that uniquely identifies a Caddy module. A
module ID is lightly structured. It consists of dot-separated
labels which form a simple hierarchy from left to right. 
The last label is the module name, and the labels before that constitutes
the namespace (or scope).
```go
type ModuleId string
```
Thus, a module ID has the form: `<namespace>.<name>`

An ID with no dot has the empty namespace, which is appropriate 
for app modules (these are "top-level" modules that Caddy core loads and runs).

Module IDs should be lowercase and use underscores (_) instead of spaces.

Examples of valid IDs:
- http
- http.handlers.file_server
- caddy.logging.encoders.json

---

* `Config` type : 

(caddy.go)

`Config` is the top (or beginning) of the Caddy configuration structure.
Caddy config is expressed natively as a JSON document.

Many parts of the config are extensible through the use of Caddy modules.
Fields which have a json.RawMessage type and which appear as dots (•••) in
the online docs can be fulfilled by modules in a certain module
namespace. The docs show which modules can be used in a given place.

Whenever a module is used, its name must be given either inline as part of 
the module or as the key to the module's value. The docs will make it clear
which to use.

Generally, all config settings are optional, as it is Caddy convention to 
have good, documented default values. If a parameter is required,
the docs should say no.

Go programs which are directly building a Config struct value should take
care of populate the JSON-encodable fields of the struct
(i.e. the fields which `json` struct tags) 
if employing the module lifecycle (e.g. `Provision` method calls).
```go
type Config struct {
	Admin   *AdminConfig `json:"admin,omitempty"`
	Logging *Logging     `json:"logging,omitempty"`

	// StorageRaw is a storage module that defines how/where Caddy
	// stores asserts (such as TLS certificates). The default storage
	// module is `caddy.storage.file_system` (the local file system),
	// and the default path
	StorageRaw json.RawMessage `json:"storage,omitempty" caddy:"namespace=caddy.storage inline_key=module"`

	// AppsRaw are the apps that Caddy will load and run.
	// The app module name is the key, 
	// and the app's config is the associated value.
	AppsRaw    ModuleMap       `json:"apps,omitempty" caddy:"namespace="`

	apps       map[string]App
	storage    certmagic.Storage
	cancelFunc context.CancelFunc
}
```

`AdminConfig`

(admin.go)

Configures Caddy's API endpoint, which is used
to manage Caddy while it is running.
```go
type AdminConfig struct {
	Disabled      bool            `json:"disabled,omitempty"`
	Listen        string          `json:"listen,omitempty"`
	EnforceOrigin bool            `json:"enforce_origin,omitempty"`
	Origins       []string        `json:"origins,omitempty"`
	Config        *ConfigSettings `json:"config,omitempty"`
	Identity      *IdentityConfig
	Remote        *RemoteAdmin
	routers       []AdminRouter
}
```

__AdminConfig fields__ :
* `Disabled` :  

If true, the admin endpoint will be completely disabled.
Note that this makes any runtime changes to the config
impossible, since the interface to do so is through the 
admin endpoint.

* `Listen` :

The address to which the admin endpoint's listener should
bind itself. Can be any single network address that can be
parsed by Caddy. Accept placeholders. Default: localhost:2019

* `EnforceOrigin` :

If true, CORS header will be emitted, and requests to the
API will be rejected if their `Host` and `Origin` headers
do not match the expected value(s). Use `origins` to 
customize which origins address is the only value allowed by
default. Enforced only on local (plaintext) endpoint

* `Origin` :

The list of allowed origins/hosts for API requests. Only needed
if accessing the admin endpoint from a host different from the
socket's network interface or if `enforce_origin` is true. If not
set, the listener address will be the default value. If set but
empty, no origins will be allowed. Enforce only on local
(plaintext) endpoint.

* `Config` :

Options pertaining to configuration management.

* `Identity` :

Options that establish this server's identity. Identity refers to
credentials which can be used to uniquely identify and authenticate
this server instance. This is required if remote administration is
enabled (but does not require remote administration to be enabled).
Default: no identity management.

* `Remote` :

Options pertaining to remote administration. By default, remote
administration is disabled. If enabled, identity management must
also be configured, as that is how the endpoint is secured.
See the neighboring "identity" object.

* `routers` :

Holds onto the routers so that we can later provision them 
if they require provisioning.

---

`Logging`

(logging.go)

`Logging` facilitates logging within Caddy. The default log is
call "default" and you can customize it. You can also define additional logs.

By default, all logs at INFO level and higher are written to 
standard error ("stderr" writer) in a human-readable format
("console" encoder if stdout is an interactive terminal, 
"json" encoder otherwise).

All defined logs accept all log entries by default, but you
can filter by level and module/logger names. A logger's name
is the same as the module's name, but module may append to 
logger names for more specificity. For example, you can
filter logs emitted only by HTTP handlers using the name
"http.handlers", because all HTTP handler module names have that prefix.

Caddy logs (expect the sink) are zero-allocation, 
so they are very high-performing in terms of memory and CPU time.
Enabling sampling can further increase throughput on extremely high-load servers.
```go
type Logging struct {
	// Sink is the destination for all unstructured logs emitted 
	// from Go's standard library logger. These logs are common
	// in dependencies that are not designed specificity for use
	// in Caddy. Because it is global, and unstructured, the sink
	// lacks most advanced features and customizations.
	Sink       *StandardLibLog `json:"sink,omitempty"`

	// Logs are your logs, keyed by an arbitrary name of your
	// choosing. The default log can be customized by defining
	// a log called "default". You can further define other logs
	// and filter what kinds of entries they accept.
	Logs       map[string]*CustomLog `json:"logs,omitempty"`

	// a list of all keys for open writers; all writers
	// that are opened to provision this logging config 
	// must have their keys added to this list so they 
	// can be closed when cleaning up
	writerKeys []string
}
```

---

(logging.go)
```go
type StandardLibLog struct {
	// The module that writes out log entries for the sink.
	WriterRaw json.RawMessage `json:"writer,omitempty" caddy:"namespace=caddy.logging.writers inline_key=output"`
	write     io.WriteCloser
}
```
StandardLibLog configures the default Go standard library
global logger in the log package. This is necessary because
module dependencies which are not built specifically for
Caddy will use the standard logger. This is also known as 
the "sink" logger.

---

(logging.go)
```go
type CustomLog struct {
	// The writer defines where log entries are emitted.
	WriteRaw   json.RawMessage `json:"writer,omitempty" caddy:"namespace=caddy.logging.writers inline_key=output"`

	// The encoder is how the log entries are formatted or encoded.
	EncoderRaw json.RawMessage `json:"encoder,omitempty" caddy:"namespace=caddy.logging.encoders inline_key=format"`

	// Level is teh minimum level to emit, and is inclusive.
	// Possible levels: DEBUG, INFO, WARN, ERROR, PANIC, and FATAL
	Level      string          `json:"level,omitempty"`

	// Sampling configures log entry sampling. If enabled,
	// only some log entries will be emitted. This is useful
	// for improving performance on extremely high-performance servers.
	Sampling   *LogSampling    `json:"sampling,omitempty"`

	// Include defines the names of loggers to emit in this
	// log. For example, to include only logs emitted by the
	// admin API, you would include "admin.api".
	Include    []string        `json:"include,omitempty"`

	// Exclude defines the names of loggers that should be
	// skipped by this log. For example, to exclude only
	// HTTP access logs, you would exclude "http.log.access".
	Exclude    []string        `json:"encoding,omitempty"`

	writerOpener WriterOpener
	writer       io.WriteCloser
	encoder      zapcore.Encoder
	levelEnable  zapcore.LevelEnabler
	core         zapcore.Core
}
```
`CustomLog` represents a custom logger configuration

By default, a log wil emit all log entries. 
Some entries will be skipped if sampling is enabled.
Further, the Include and Exclude parameters define which 
loggers (by name) are allowed or rejected from emitting in this log.
If both Include and Exclude are populated, their values must be mutually
are populated, all logs are emitted.

---

(logging.go)
```go
type LogSampling struct {
	// The window over which to conduct sampling.
	Interval   time.Duration `json:"interval,omitempty"`

	// Log this many entries within a given level and 
	// message for each interval.
	First      int           `json:first,omitempty"`

	// If more entries with the same level and message 
	// are seen during the same interval, keep one in 
	// this many entries until the end of the interval.
	Thereafter int           `json:thereafter,omitempty"`
}
```
`LogSampling` configures log entry sampling.

---

(logging.go)
```go
type WriterOpener interface {
	fmt.Stringer

	// WriterKey is a string that uniquely identifies this
	// writer configuration. It is not shown to humans.
	WriterKey() string

	// OpenWriter opens a log for writing. The writer
	// should be safe for concurrent use but need not 
	// be synchronous.
	OpenWriter() (io.WriteCloser, error)
}
```
`WriterOpener` is a module that can open a log writer.
It can return a human-readable string representation of
itself so that operators can understand where the logs are going.

---

(module.go)
```go
type ModuleMap map[string]json.RawMessage
```
`ModuleMap` is a map that can contain multiple modules.
where the map key is the module's name. (The namespace
is usually read from an associated field's struct tag.)
module map, the name does not have to be given in the json.RawMessage.

---

```go
type App interface {
	Start() error
	Stop() error
}
```
`App` is a thing that Caddy runs.

---


__Context Functions__ : 
* `func NewContext(ctx Context) (Context, context.CancelFunc)`

`NewContext` provides a new context derived from the given
context ctx. Normally, you will not need to call this
function unless you are loading modules which have a
module was provisioned with. Be sure to call the cancel
func when the context is to be cleaned up so that
modules which are loaded will be properly unloaded.
See standard library context package's documentation.



