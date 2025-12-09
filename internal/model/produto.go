package model

import "time"

// Produto representa a entidade de domínio de um produto
// Esta é a entidade interna usada para persistência e lógica de negócio
type Produto struct {
	ID        int        `json:"id" bson:"id"`
	Nome      string     `json:"nome" bson:"nome"`
	Preco     float64    `json:"preco" bson:"preco"`
	Descricao string     `json:"descricao" bson:"descricao"`
	CreatedAt time.Time  `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" bson:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" bson:"deleted_at,omitempty"` // Soft delete
}

// IsDeleted verifica se o produto foi deletado (soft delete)
func (p *Produto) IsDeleted() bool {
	return p.DeletedAt != nil && !p.DeletedAt.IsZero()
}

// SoftDelete marca o produto como deletado (soft delete)
func (p *Produto) SoftDelete() {
	now := time.Now()
	p.DeletedAt = &now
	p.UpdatedAt = now
}

// Restore restaura um produto deletado (soft delete)
func (p *Produto) Restore() {
	p.DeletedAt = nil
	p.UpdatedAt = time.Now()
}

// BeforeCreate inicializa os timestamps antes de criar
func (p *Produto) BeforeCreate() {
	now := time.Now()
	if p.CreatedAt.IsZero() {
		p.CreatedAt = now
	}
	if p.UpdatedAt.IsZero() {
		p.UpdatedAt = now
	}
}

// BeforeUpdate atualiza o timestamp de atualização
func (p *Produto) BeforeUpdate() {
	p.UpdatedAt = time.Now()
}
