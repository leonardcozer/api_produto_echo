package cache

import (
	"encoding/json"
	"api-go-arquitetura/internal/model"
)

// EncodeProduto codifica um produto para JSON
func EncodeProduto(produto model.Produto) ([]byte, error) {
	return json.Marshal(produto)
}

// DecodeProduto decodifica um produto de JSON
func DecodeProduto(data []byte) (model.Produto, error) {
	var produto model.Produto
	err := json.Unmarshal(data, &produto)
	return produto, err
}

// EncodeProdutos codifica uma lista de produtos para JSON
func EncodeProdutos(produtos []model.Produto) ([]byte, error) {
	return json.Marshal(produtos)
}

// DecodeProdutos decodifica uma lista de produtos de JSON
func DecodeProdutos(data []byte) ([]model.Produto, error) {
	var produtos []model.Produto
	err := json.Unmarshal(data, &produtos)
	return produtos, err
}

// Encode codifica qualquer valor para JSON (genérico)
func Encode(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

// Decode decodifica qualquer valor de JSON (genérico)
func Decode(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

