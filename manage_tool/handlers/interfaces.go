package handlers

import (
	"strikepad-manage-tool/repository/ri"
	"strikepad-manage-tool/handlers/hi"
)

// Repository Interfaces - Using ri package

// MakerRepositoryInterface はメーカーリポジトリのインターフェース
type MakerRepositoryInterface = ri.MakerRepositoryInterface

// UserRepositoryInterface はユーザーリポジトリのインターフェース
type UserRepositoryInterface = ri.UserRepositoryInterface

// AdminRepositoryInterface は管理者リポジトリのインターフェース
type AdminRepositoryInterface = ri.AdminRepositoryInterface

// CoreRepositoryInterface はコアリポジトリのインターフェース
type CoreRepositoryInterface = ri.CoreRepositoryInterface

// CoverRepositoryInterface はカバーリポジトリのインターフェース
type CoverRepositoryInterface = ri.CoverRepositoryInterface

// Handler Interfaces - Using hi package

// AuthHandlerInterface は認証ハンドラーのインターフェース
type AuthHandlerInterface = hi.AuthHandlerInterface

// AdminHandlerInterface は管理者ハンドラーのインターフェース
type AdminHandlerInterface = hi.AdminHandlerInterface

// MakerHandlerInterface はメーカーハンドラーのインターフェース
type MakerHandlerInterface = hi.MakerHandlerInterface

// CoreHandlerInterface はコアハンドラーのインターフェース
type CoreHandlerInterface = hi.CoreHandlerInterface

// CoverHandlerInterface はカバーハンドラーのインターフェース
type CoverHandlerInterface = hi.CoverHandlerInterface
