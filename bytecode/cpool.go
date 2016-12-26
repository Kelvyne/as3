package bytecode

import "fmt"

func typenameString(c *CpoolInfo, info MultinameInfo) string {
	var str string
	str += c.Strings[info.Name]
	str += "<"
	for i, p := range info.Params {
		if i > 0 {
			str += ", "
		}
		str += c.MultinameString(p)
	}
	str += ">"
	return str
}

// MultinameString converts a multiname to a string
func (c *CpoolInfo) MultinameString(m uint32) string {
	info := c.Multinames[m]
	switch info.Kind {
	case MultinameKindQName, MultinameKindQNameA:
		return fmt.Sprintf("[%v].%v",
			c.NamespaceString(info.Namespace),
			c.Strings[info.Name])
	case MultinameKindRTQName, MultinameKindRTQNameA:
		return fmt.Sprintf("[*].%v", c.Strings[info.Name])
	case MultinameKindRTQNameL, MultinameKindRTQNameLA:
		return fmt.Sprint("[*].[*]")
	case MultinameKindTypename:
		return typenameString(c, info)
	default:
		return fmt.Sprint(c.Strings[info.Name])
	}
}

// NamespaceString converts a namespace to a string
func (c *CpoolInfo) NamespaceString(n uint32) string {
	info := c.Namespaces[n]
	return c.Strings[info.Name]
}
