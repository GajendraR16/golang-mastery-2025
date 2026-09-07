package main

type InMemoryStore struct {
	tasks map[int]*Task
}

func NewInMemoryStoreTask() *InMemoryStore {
	return &InMemoryStore{
		tasks: make(map[int]*Task),
	}
}

func (im *InMemoryStore) Save(task *Task) error {
	im.tasks[task.ID] = task
	return nil
}

func (im *InMemoryStore) Get(id int) (*Task, error) {
	if task, ok := im.tasks[id]; ok {
		return task, nil
	}

	return nil, TaskNotFoundError{ID: id}
}

func (im *InMemoryStore) Delete(id int) error {
	if _, ok := im.tasks[id]; ok {
		delete(im.tasks, id)
		return nil
	}

	return TaskNotFoundError{ID: id}
}

func (im *InMemoryStore) List() ([]*Task, error) {
	tasks := []*Task{}

	for _, task := range im.tasks {
		tasks = append(tasks, task)
	}

	return tasks, nil
}
