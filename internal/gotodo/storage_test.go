package gotodo

type TestStorage struct {
	items TodoList
}

func (me *TestStorage) Create(todo *Todo) error {
	me.items = append(me.items, todo)
	return nil
}

func (me *TestStorage) List() (TodoList, error) {
	return me.items, nil
}

func (me *TestStorage) Get(todoID int) (*Todo, error) {
	return me.items[todoID], nil
}

func (me *TestStorage) Update(todoID int, todo *Todo) error {
	me.items[todoID] = todo
	return nil
}

func (me *TestStorage) Delete(todoID int) error {
	todos := make(TodoList, 0)
	left := me.items[:todoID]
	right := me.items[todoID+1:]
	todos = append(todos, left...)
	todos = append(todos, right...)
	me.items = todos
	return nil
}
