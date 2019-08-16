package script

type Component interface {
    GetOwner() interface{}
}

type ComponentType struct {
    owner interface{}
}

func (o *ComponentType) GetOwner() interface{} {
    return o.owner
}

func MakeComponentType(owner interface{}) ComponentType {
    return ComponentType{owner: owner}
}
