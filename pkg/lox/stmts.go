package lox

// Visitor pattern for Stmts
type Stmt interface {
	Accept(StmtVisitor) error
}

type StmtVisitor interface {
	visitBlockStmt(BlockStmt) error
	visitExprStmt(ExprStmt) error
	visitFuncStmt(FuncStmt) error
	visitIfStmt(IfStmt) error
	visitPrintStmt(PrintStmt) error
	visitReturnStmt(ReturnStmt) error
	visitVarStmt(VarStmt) error
	visitWhileStmt(WhileStmt) error
}

type BlockStmt struct {
	statements []Stmt
}

func (b BlockStmt) Accept(visitor StmtVisitor) error {
	return visitor.visitBlockStmt(b)
}

type ExprStmt struct {
	expression Expr
}

func (e ExprStmt) Accept(visitor StmtVisitor) error {
	return visitor.visitExprStmt(e)
}

type FuncStmt struct {
	name    Token
	params  []Token
	body    []Stmt
	closure *Environment
}

func (f FuncStmt) Accept(visitor StmtVisitor) error {
	return visitor.visitFuncStmt(f)
}

type IfStmt struct {
	condition Expr
	branch    Stmt
	elseStmt  Stmt
}

func (i IfStmt) Accept(visitor StmtVisitor) error {
	return visitor.visitIfStmt(i)
}

type PrintStmt struct {
	expression Expr
}

func (p PrintStmt) Accept(visitor StmtVisitor) error {
	return visitor.visitPrintStmt(p)
}

type ReturnStmt struct {
	keyword Token
	value   Expr
}

func (r ReturnStmt) Accept(visitor StmtVisitor) error {
	return visitor.visitReturnStmt(r)
}

type VarStmt struct {
	name        Token
	initializer Expr
}

func (v VarStmt) Accept(visitor StmtVisitor) error {
	return visitor.visitVarStmt(v)
}

type WhileStmt struct {
	condition Expr
	body      Stmt
}

func (w WhileStmt) Accept(visitor StmtVisitor) error {
	return visitor.visitWhileStmt(w)
}
