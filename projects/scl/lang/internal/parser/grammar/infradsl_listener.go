// Code generated from grammar/InfraDSL.g4 by ANTLR 4.13.1. DO NOT EDIT.

package parser // InfraDSL

import "github.com/antlr4-go/antlr/v4"

// InfraDSLListener is a complete listener for a parse tree produced by InfraDSLParser.
type InfraDSLListener interface {
	antlr.ParseTreeListener

	// EnterProgram is called when entering the program production.
	EnterProgram(c *ProgramContext)

	// EnterImportStatement is called when entering the importStatement production.
	EnterImportStatement(c *ImportStatementContext)

	// EnterStatement is called when entering the statement production.
	EnterStatement(c *StatementContext)

	// EnterAssignment is called when entering the assignment production.
	EnterAssignment(c *AssignmentContext)

	// EnterDeclaration is called when entering the declaration production.
	EnterDeclaration(c *DeclarationContext)

	// EnterType is called when entering the type production.
	EnterType(c *TypeContext)

	// EnterFunctionDeclaration is called when entering the functionDeclaration production.
	EnterFunctionDeclaration(c *FunctionDeclarationContext)

	// EnterParameterList is called when entering the parameterList production.
	EnterParameterList(c *ParameterListContext)

	// EnterParameter is called when entering the parameter production.
	EnterParameter(c *ParameterContext)

	// EnterBlock is called when entering the block production.
	EnterBlock(c *BlockContext)

	// EnterExpressionStatement is called when entering the expressionStatement production.
	EnterExpressionStatement(c *ExpressionStatementContext)

	// EnterFunctionCallExpr is called when entering the FunctionCallExpr production.
	EnterFunctionCallExpr(c *FunctionCallExprContext)

	// EnterMulDivModExpr is called when entering the MulDivModExpr production.
	EnterMulDivModExpr(c *MulDivModExprContext)

	// EnterComparisonExpr is called when entering the ComparisonExpr production.
	EnterComparisonExpr(c *ComparisonExprContext)

	// EnterPrimaryExpr is called when entering the PrimaryExpr production.
	EnterPrimaryExpr(c *PrimaryExprContext)

	// EnterNotExpr is called when entering the NotExpr production.
	EnterNotExpr(c *NotExprContext)

	// EnterMemberAccessExpr is called when entering the MemberAccessExpr production.
	EnterMemberAccessExpr(c *MemberAccessExprContext)

	// EnterAddSubExpr is called when entering the AddSubExpr production.
	EnterAddSubExpr(c *AddSubExprContext)

	// EnterLogicalExpr is called when entering the LogicalExpr production.
	EnterLogicalExpr(c *LogicalExprContext)

	// EnterQualifiedName is called when entering the qualifiedName production.
	EnterQualifiedName(c *QualifiedNameContext)

	// EnterArgumentList is called when entering the argumentList production.
	EnterArgumentList(c *ArgumentListContext)

	// EnterIfStatement is called when entering the ifStatement production.
	EnterIfStatement(c *IfStatementContext)

	// EnterPrimary is called when entering the primary production.
	EnterPrimary(c *PrimaryContext)

	// EnterArray is called when entering the array production.
	EnterArray(c *ArrayContext)

	// ExitProgram is called when exiting the program production.
	ExitProgram(c *ProgramContext)

	// ExitImportStatement is called when exiting the importStatement production.
	ExitImportStatement(c *ImportStatementContext)

	// ExitStatement is called when exiting the statement production.
	ExitStatement(c *StatementContext)

	// ExitAssignment is called when exiting the assignment production.
	ExitAssignment(c *AssignmentContext)

	// ExitDeclaration is called when exiting the declaration production.
	ExitDeclaration(c *DeclarationContext)

	// ExitType is called when exiting the type production.
	ExitType(c *TypeContext)

	// ExitFunctionDeclaration is called when exiting the functionDeclaration production.
	ExitFunctionDeclaration(c *FunctionDeclarationContext)

	// ExitParameterList is called when exiting the parameterList production.
	ExitParameterList(c *ParameterListContext)

	// ExitParameter is called when exiting the parameter production.
	ExitParameter(c *ParameterContext)

	// ExitBlock is called when exiting the block production.
	ExitBlock(c *BlockContext)

	// ExitExpressionStatement is called when exiting the expressionStatement production.
	ExitExpressionStatement(c *ExpressionStatementContext)

	// ExitFunctionCallExpr is called when exiting the FunctionCallExpr production.
	ExitFunctionCallExpr(c *FunctionCallExprContext)

	// ExitMulDivModExpr is called when exiting the MulDivModExpr production.
	ExitMulDivModExpr(c *MulDivModExprContext)

	// ExitComparisonExpr is called when exiting the ComparisonExpr production.
	ExitComparisonExpr(c *ComparisonExprContext)

	// ExitPrimaryExpr is called when exiting the PrimaryExpr production.
	ExitPrimaryExpr(c *PrimaryExprContext)

	// ExitNotExpr is called when exiting the NotExpr production.
	ExitNotExpr(c *NotExprContext)

	// ExitMemberAccessExpr is called when exiting the MemberAccessExpr production.
	ExitMemberAccessExpr(c *MemberAccessExprContext)

	// ExitAddSubExpr is called when exiting the AddSubExpr production.
	ExitAddSubExpr(c *AddSubExprContext)

	// ExitLogicalExpr is called when exiting the LogicalExpr production.
	ExitLogicalExpr(c *LogicalExprContext)

	// ExitQualifiedName is called when exiting the qualifiedName production.
	ExitQualifiedName(c *QualifiedNameContext)

	// ExitArgumentList is called when exiting the argumentList production.
	ExitArgumentList(c *ArgumentListContext)

	// ExitIfStatement is called when exiting the ifStatement production.
	ExitIfStatement(c *IfStatementContext)

	// ExitPrimary is called when exiting the primary production.
	ExitPrimary(c *PrimaryContext)

	// ExitArray is called when exiting the array production.
	ExitArray(c *ArrayContext)
}
