package view

import (
	"html/template"
	"net/http"
)

// extend safely reads the extend list.
func (c *Info) extend() template.FuncMap {
	extendMutex.RLock()
	list := c.extendList
	extendMutex.RUnlock()

	return list
}

// modify safely reads the modify list.
func (c *Info) modify() []ModifyFunc {
	// Get the setter collection
	modifyMutex.RLock()
	list := c.modifyList
	modifyMutex.RUnlock()

	return list
}

// SetTemplates will set the root and child templates.
func (c *Info) SetTemplates(rootTemp string, childTemps []string) {
	templateCacheMutex.Lock()
	c.templateCollection = make(map[string]*template.Template)
	templateCacheMutex.Unlock()

	c.rootTemplate = rootTemp
	c.childTemplates = childTemps
}

// ModifyFunc can modify the view before rendering.
type ModifyFunc func(http.ResponseWriter, *http.Request, *Info)

// SetModifiers will set the modifiers for the View that run
// before rendering.
func (c *Info) SetModifiers(fn ...ModifyFunc) {
	modifyMutex.Lock()
	c.modifyList = fn
	modifyMutex.Unlock()
}

// SetFuncMaps will combine all template.FuncMaps into one map and then set the
// them for each template.
// If a func already exists, it is rewritten without a warning.
func (c *Info) SetFuncMaps(fms ...template.FuncMap) {
	// Final FuncMap
	fm := make(template.FuncMap)

	// Loop through the maps
	for _, m := range fms {
		// Loop through each key and value
		for k, v := range m {
			fm[k] = v
		}
	}

	// Load the plugins
	extendMutex.Lock()
	c.extendList = fm
	extendMutex.Unlock()
}
