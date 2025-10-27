// Code generated from InfraDSL.g4 by ANTLR 4.13.1. DO NOT EDIT.

package parser // InfraDSL

import "github.com/antlr4-go/antlr/v4"

// BaseInfraDSLListener is a complete listener for a parse tree produced by InfraDSLParser.
type BaseInfraDSLListener struct{}

var _ InfraDSLListener = &BaseInfraDSLListener{}

// VisitTerminal is called when a terminal node is visited.
func (s *BaseInfraDSLListener) VisitTerminal(node antlr.TerminalNode) {}

// VisitErrorNode is called when an error node is visited.
func (s *BaseInfraDSLListener) VisitErrorNode(node antlr.ErrorNode) {}

// EnterEveryRule is called when any rule is entered.
func (s *BaseInfraDSLListener) EnterEveryRule(ctx antlr.ParserRuleContext) {}

// ExitEveryRule is called when any rule is exited.
func (s *BaseInfraDSLListener) ExitEveryRule(ctx antlr.ParserRuleContext) {}

// EnterProgram is called when production program is entered.
func (s *BaseInfraDSLListener) EnterProgram(ctx *ProgramContext) {}

// ExitProgram is called when production program is exited.
func (s *BaseInfraDSLListener) ExitProgram(ctx *ProgramContext) {}

// EnterImportStatement is called when production importStatement is entered.
func (s *BaseInfraDSLListener) EnterImportStatement(ctx *ImportStatementContext) {}

// ExitImportStatement is called when production importStatement is exited.
func (s *BaseInfraDSLListener) ExitImportStatement(ctx *ImportStatementContext) {}

// EnterStatement is called when production statement is entered.
func (s *BaseInfraDSLListener) EnterStatement(ctx *StatementContext) {}

// ExitStatement is called when production statement is exited.
func (s *BaseInfraDSLListener) ExitStatement(ctx *StatementContext) {}

// EnterAssignment is called when production assignment is entered.
func (s *BaseInfraDSLListener) EnterAssignment(ctx *AssignmentContext) {}

// ExitAssignment is called when production assignment is exited.
func (s *BaseInfraDSLListener) ExitAssignment(ctx *AssignmentContext) {}

// EnterDeclaration is called when production declaration is entered.
func (s *BaseInfraDSLListener) EnterDeclaration(ctx *DeclarationContext) {}

// ExitDeclaration is called when production declaration is exited.
func (s *BaseInfraDSLListener) ExitDeclaration(ctx *DeclarationContext) {}

// EnterType is called when production type is entered.
func (s *BaseInfraDSLListener) EnterType(ctx *TypeContext) {}

// ExitType is called when production type is exited.
func (s *BaseInfraDSLListener) ExitType(ctx *TypeContext) {}

// EnterFunctionDeclaration is called when production functionDeclaration is entered.
func (s *BaseInfraDSLListener) EnterFunctionDeclaration(ctx *FunctionDeclarationContext) {}

// ExitFunctionDeclaration is called when production functionDeclaration is exited.
func (s *BaseInfraDSLListener) ExitFunctionDeclaration(ctx *FunctionDeclarationContext) {}

// EnterParameterList is called when production parameterList is entered.
func (s *BaseInfraDSLListener) EnterParameterList(ctx *ParameterListContext) {}

// ExitParameterList is called when production parameterList is exited.
func (s *BaseInfraDSLListener) ExitParameterList(ctx *ParameterListContext) {}

// EnterParameter is called when production parameter is entered.
func (s *BaseInfraDSLListener) EnterParameter(ctx *ParameterContext) {}

// ExitParameter is called when production parameter is exited.
func (s *BaseInfraDSLListener) ExitParameter(ctx *ParameterContext) {}

// EnterBlock is called when production block is entered.
func (s *BaseInfraDSLListener) EnterBlock(ctx *BlockContext) {}

// ExitBlock is called when production block is exited.
func (s *BaseInfraDSLListener) ExitBlock(ctx *BlockContext) {}

// EnterExpressionStatement is called when production expressionStatement is entered.
func (s *BaseInfraDSLListener) EnterExpressionStatement(ctx *ExpressionStatementContext) {}

// ExitExpressionStatement is called when production expressionStatement is exited.
func (s *BaseInfraDSLListener) ExitExpressionStatement(ctx *ExpressionStatementContext) {}

// EnterFunctionCallExpr is called when production FunctionCallExpr is entered.
func (s *BaseInfraDSLListener) EnterFunctionCallExpr(ctx *FunctionCallExprContext) {}

// ExitFunctionCallExpr is called when production FunctionCallExpr is exited.
func (s *BaseInfraDSLListener) ExitFunctionCallExpr(ctx *FunctionCallExprContext) {}

// EnterMulDivModExpr is called when production MulDivModExpr is entered.
func (s *BaseInfraDSLListener) EnterMulDivModExpr(ctx *MulDivModExprContext) {}

// ExitMulDivModExpr is called when production MulDivModExpr is exited.
func (s *BaseInfraDSLListener) ExitMulDivModExpr(ctx *MulDivModExprContext) {}

// EnterComparisonExpr is called when production ComparisonExpr is entered.
func (s *BaseInfraDSLListener) EnterComparisonExpr(ctx *ComparisonExprContext) {}

// ExitComparisonExpr is called when production ComparisonExpr is exited.
func (s *BaseInfraDSLListener) ExitComparisonExpr(ctx *ComparisonExprContext) {}

// EnterPrimaryExpr is called when production PrimaryExpr is entered.
func (s *BaseInfraDSLListener) EnterPrimaryExpr(ctx *PrimaryExprContext) {}

// ExitPrimaryExpr is called when production PrimaryExpr is exited.
func (s *BaseInfraDSLListener) ExitPrimaryExpr(ctx *PrimaryExprContext) {}

// EnterNotExpr is called when production NotExpr is entered.
func (s *BaseInfraDSLListener) EnterNotExpr(ctx *NotExprContext) {}

// ExitNotExpr is called when production NotExpr is exited.
func (s *BaseInfraDSLListener) ExitNotExpr(ctx *NotExprContext) {}

// EnterMemberAccessExpr is called when production MemberAccessExpr is entered.
func (s *BaseInfraDSLListener) EnterMemberAccessExpr(ctx *MemberAccessExprContext) {}

// ExitMemberAccessExpr is called when production MemberAccessExpr is exited.
func (s *BaseInfraDSLListener) ExitMemberAccessExpr(ctx *MemberAccessExprContext) {}

// EnterAddSubExpr is called when production AddSubExpr is entered.
func (s *BaseInfraDSLListener) EnterAddSubExpr(ctx *AddSubExprContext) {}

// ExitAddSubExpr is called when production AddSubExpr is exited.
func (s *BaseInfraDSLListener) ExitAddSubExpr(ctx *AddSubExprContext) {}

// EnterLogicalExpr is called when production LogicalExpr is entered.
func (s *BaseInfraDSLListener) EnterLogicalExpr(ctx *LogicalExprContext) {}

// ExitLogicalExpr is called when production LogicalExpr is exited.
func (s *BaseInfraDSLListener) ExitLogicalExpr(ctx *LogicalExprContext) {}

// EnterQualifiedName is called when production qualifiedName is entered.
func (s *BaseInfraDSLListener) EnterQualifiedName(ctx *QualifiedNameContext) {}

// ExitQualifiedName is called when production qualifiedName is exited.
func (s *BaseInfraDSLListener) ExitQualifiedName(ctx *QualifiedNameContext) {}

// EnterArgumentList is called when production argumentList is entered.
func (s *BaseInfraDSLListener) EnterArgumentList(ctx *ArgumentListContext) {}

// ExitArgumentList is called when production argumentList is exited.
func (s *BaseInfraDSLListener) ExitArgumentList(ctx *ArgumentListContext) {}

// EnterIfStatement is called when production ifStatement is entered.
func (s *BaseInfraDSLListener) EnterIfStatement(ctx *IfStatementContext) {}

// ExitIfStatement is called when production ifStatement is exited.
func (s *BaseInfraDSLListener) ExitIfStatement(ctx *IfStatementContext) {}

// EnterPrimary is called when production primary is entered.
func (s *BaseInfraDSLListener) EnterPrimary(ctx *PrimaryContext) {}

// ExitPrimary is called when production primary is exited.
func (s *BaseInfraDSLListener) ExitPrimary(ctx *PrimaryContext) {}

// EnterArray is called when production array is entered.
func (s *BaseInfraDSLListener) EnterArray(ctx *ArrayContext) {}

// ExitArray is called when production array is exited.
func (s *BaseInfraDSLListener) ExitArray(ctx *ArrayContext) {}
