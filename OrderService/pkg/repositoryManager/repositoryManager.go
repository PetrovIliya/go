package repositoryManager

import (
	"OrderService/pkg/repository"
)

type RepositoryManager struct {
	orderRepository repository.OrderRepository
}

func Create(orderRepository repository.OrderRepository) RepositoryManager {
	var repositoryManager RepositoryManager
	repositoryManager.orderRepository = orderRepository
	return repositoryManager
}

func (repositoryManager RepositoryManager) GetOrderRepository() repository.OrderRepository {
	return repositoryManager.orderRepository
}