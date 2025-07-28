package handlers


func (h *UserHandler) DeleteUserbyID(id int) error {
	return h.repo.DeleteUserByID(id)
}

func (h *bookHandler) DeleteBookByID(id int) error {
	return h.repo.DeleteBookByID(id)
}