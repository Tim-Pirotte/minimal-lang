package plugins

import (
	"reflect"
)

type registry struct {
	plugins      plugins
}

type plugins struct {
    mapping   map[reflect.Type]Plugin
}

type Plugin interface {
    Init(*plugins)
}

func NewRegistry() registry {
    return registry{plugins{map[reflect.Type]Plugin{}}}
}

func (r *registry) Add(plugin Plugin) bool {
    t := reflect.TypeOf(plugin)

    if t.Kind() != reflect.Pointer {
        panic("plugin is not a pointer")
    }

    if _, ok := r.plugins.mapping[t]; ok {
        return false
    }

    r.plugins.mapping[t] = plugin

    return true
}

func (r *registry) Setup() {
    for _, plugin := range r.plugins.mapping {
        plugin.Init(&r.plugins)
    }
}

func Get[T any](p *plugins) (T, bool) {
    t := reflect.TypeFor[T]()

    if t.Kind() != reflect.Pointer {
        panic("T should be a pointer")
    }

    if p, ok := p.mapping[t]; ok {
        return p.(T), true
    } else {
        var empty T

        return empty, false
    }
}
