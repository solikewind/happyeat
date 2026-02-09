### 1. "获取单个菜单"

1. route definition

- Url: /central/v1/menu/:id
- Method: GET
- Request: `GetMenuReq`
- Response: `GetMenuReply`

2. request definition



```golang
type GetMenuReq struct {
	Id uint64 `json:"id"`
}
```


3. response definition



```golang
type GetMenuReply struct {
	Menu Menu `json:"menu"`
}

type Menu struct {
	Id uint64 `json:"id"` // 菜单id
	Name string `json:"name"` // 菜单名称
	Image string `json:"image"` // 菜单图片
	Price float64 `json:"price"` // 菜单价格
	Category string `json:"category"` // 菜单分类
	Spec string `json:"spec"` // 菜单规格
	Description string `json:"description"` // 菜单描述
	Create_at int64 `json:"create_at"` // 创建时间
}
```

### 2. "更新菜单"

1. route definition

- Url: /central/v1/menu/:id
- Method: PUT
- Request: `UpdateMenuReq`
- Response: `UpdateMenuReply`

2. request definition



```golang
type UpdateMenuReq struct {
	Menu Menu `json:"menu"`
}

type Menu struct {
	Id uint64 `json:"id"` // 菜单id
	Name string `json:"name"` // 菜单名称
	Image string `json:"image"` // 菜单图片
	Price float64 `json:"price"` // 菜单价格
	Category string `json:"category"` // 菜单分类
	Spec string `json:"spec"` // 菜单规格
	Description string `json:"description"` // 菜单描述
	Create_at int64 `json:"create_at"` // 创建时间
}
```


3. response definition



```golang
type UpdateMenuReply struct {
}
```

### 3. "删除菜单"

1. route definition

- Url: /central/v1/menu/:id
- Method: DELETE
- Request: `DeleteMenuReq`
- Response: `DeleteMenuReply`

2. request definition



```golang
type DeleteMenuReq struct {
	Id uint64 `json:"id"`
}
```


3. response definition



```golang
type DeleteMenuReply struct {
}
```

### 4. "列出菜单"

1. route definition

- Url: /central/v1/menus
- Method: GET
- Request: `ListMenusReq`
- Response: `ListMenusReply`

2. request definition



```golang
type ListMenusReq struct {
	Current uint64 `json:"current"`
	PageSize uint64 `json:"pageSize"`
	Name string `json:"name,optional"`
	Category string `json:"category,optional"`
	Spec string `json:"spec,optional"`
}
```


3. response definition



```golang
type ListMenusReply struct {
	Menus []Menu `json:"menus"`
	Total uint64 `json:"total"`
}
```

### 5. "创建菜单"

1. route definition

- Url: /central/v1/menus
- Method: POST
- Request: `CreateMenuReq`
- Response: `CreateMenuReply`

2. request definition



```golang
type CreateMenuReq struct {
	Menu Menu `json:"menu"`
}

type Menu struct {
	Id uint64 `json:"id"` // 菜单id
	Name string `json:"name"` // 菜单名称
	Image string `json:"image"` // 菜单图片
	Price float64 `json:"price"` // 菜单价格
	Category string `json:"category"` // 菜单分类
	Spec string `json:"spec"` // 菜单规格
	Description string `json:"description"` // 菜单描述
	Create_at int64 `json:"create_at"` // 创建时间
}
```


3. response definition



```golang
type CreateMenuReply struct {
}
```

### 6. "列出菜单种类"

1. route definition

- Url: /central/v1/menu/categories
- Method: GET
- Request: `ListMenusCategoriesReq`
- Response: `ListMenusCategoriesReply`

2. request definition



```golang
type ListMenusCategoriesReq struct {
	Current uint64 `json:"current"`
	PageSize uint64 `json:"pageSize"`
	Name string `json:"name,optional"`
}
```


3. response definition



```golang
type ListMenusCategoriesReply struct {
	Categories []MenuCategory `json:"categories"`
	Total uint64 `json:"total"`
}
```

### 7. "创建菜单种类"

1. route definition

- Url: /central/v1/menu/category
- Method: POST
- Request: `CreateMenuCategoryReq`
- Response: `CreateMenuCategoryReply`

2. request definition



```golang
type CreateMenuCategoryReq struct {
	MenuCategory MenuCategory `json:"category"`
}

type MenuCategory struct {
	Id uint64 `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
}
```


3. response definition



```golang
type CreateMenuCategoryReply struct {
}
```

### 8. "获取菜单种类"

1. route definition

- Url: /central/v1/menu/category/:id
- Method: GET
- Request: `GetMenuCategoryReq`
- Response: `GetMenuCategoryReply`

2. request definition



```golang
type GetMenuCategoryReq struct {
	Id uint64 `json:"id"`
}
```


3. response definition



```golang
type GetMenuCategoryReply struct {
	MenuCategory MenuCategory `json:"category"`
}

type MenuCategory struct {
	Id uint64 `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
}
```

### 9. "更新菜单种类"

1. route definition

- Url: /central/v1/menu/category/:id
- Method: PUT
- Request: `UpdateMenuCategoryReq`
- Response: `UpdateMenuCategoryReply`

2. request definition



```golang
type UpdateMenuCategoryReq struct {
	MenuCategory MenuCategory `json:"category"`
}

type MenuCategory struct {
	Id uint64 `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
}
```


3. response definition



```golang
type UpdateMenuCategoryReply struct {
}
```

### 10. "删除菜单种类"

1. route definition

- Url: /central/v1/menu/category/:id
- Method: DELETE
- Request: `DeleteMenuCategoryReq`
- Response: `DeleteMenuCategoryReply`

2. request definition



```golang
type DeleteMenuCategoryReq struct {
	Id uint64 `json:"id"`
}
```


3. response definition



```golang
type DeleteMenuCategoryReply struct {
}
```

