package lox

// Visitor pattern for Stmts
type Stmt interface {
	Accept(StmtVisitor) error
}

type StmtVisitor interface {
	visitBlockStmt(BlockStmt) error
	visitPrintStmt(PrintStmt) error
	visitExprStmt(ExprStmt) error
	visitVarStmt(VarStmt) error
}

type BlockStmt struct {
	statements []Stmt
}

func (b BlockStmt) Accept(visitor StmtVisitor) error {
	return visitor.visitBlockStmt(b)
}

type PrintStmt struct {
	expression Expr
}

func (p PrintStmt) Accept(visitor StmtVisitor) error {
	return visitor.visitPrintStmt(p)
}

type ExprStmt struct {
	expression Expr
}

func (e ExprStmt) Accept(visitor StmtVisitor) error {
	return visitor.visitExprStmt(e)
}

type VarStmt struct {
	name        Token
	initializer Expr
}

func (v VarStmt) Accept(visitor StmtVisitor) error {
	return visitor.visitVarStmt(v)
}
