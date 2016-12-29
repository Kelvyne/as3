package as3

// GetClassByName finds a class by its name.
// If multiple classes have the same name (but not the same namespace)
// the first occurence is returned
func (f AbcFile) GetClassByName(name string) (Class, bool) {
	for _, c := range f.Classes {
		if c.Name == name {
			return c, true
		}
	}
	return Class{}, false
}
