// Code generated from huloParser.g4 by ANTLR 4.13.2. DO NOT EDIT.

package generated // huloParser
import "github.com/antlr4-go/antlr/v4"

// A complete Visitor for a parse tree produced by huloParser.
type huloParserVisitor interface {
	antlr.ParseTreeVisitor

	// Visit a parse tree produced by huloParser#file.
	VisitFile(ctx *FileContext) interface{}

	// Visit a parse tree produced by huloParser#block.
	VisitBlock(ctx *BlockContext) interface{}

	// Visit a parse tree produced by huloParser#comment.
	VisitComment(ctx *CommentContext) interface{}

	// Visit a parse tree produced by huloParser#statement.
	VisitStatement(ctx *StatementContext) interface{}

	// Visit a parse tree produced by huloParser#assignStatement.
	VisitAssignStatement(ctx *AssignStatementContext) interface{}

	// Visit a parse tree produced by huloParser#lambdaAssignStatement.
	VisitLambdaAssignStatement(ctx *LambdaAssignStatementContext) interface{}

	// Visit a parse tree produced by huloParser#variableNullableExpressions.
	VisitVariableNullableExpressions(ctx *VariableNullableExpressionsContext) interface{}

	// Visit a parse tree produced by huloParser#variableNullableExpression.
	VisitVariableNullableExpression(ctx *VariableNullableExpressionContext) interface{}

	// Visit a parse tree produced by huloParser#variableNames.
	VisitVariableNames(ctx *VariableNamesContext) interface{}

	// Visit a parse tree produced by huloParser#variableNameList.
	VisitVariableNameList(ctx *VariableNameListContext) interface{}

	// Visit a parse tree produced by huloParser#conditionalExpression.
	VisitConditionalExpression(ctx *ConditionalExpressionContext) interface{}

	// Visit a parse tree produced by huloParser#conditionalBoolExpression.
	VisitConditionalBoolExpression(ctx *ConditionalBoolExpressionContext) interface{}

	// Visit a parse tree produced by huloParser#logicalExpression.
	VisitLogicalExpression(ctx *LogicalExpressionContext) interface{}

	// Visit a parse tree produced by huloParser#shiftExpression.
	VisitShiftExpression(ctx *ShiftExpressionContext) interface{}

	// Visit a parse tree produced by huloParser#addSubExpression.
	VisitAddSubExpression(ctx *AddSubExpressionContext) interface{}

	// Visit a parse tree produced by huloParser#mulDivExpression.
	VisitMulDivExpression(ctx *MulDivExpressionContext) interface{}

	// Visit a parse tree produced by huloParser#incDecExpression.
	VisitIncDecExpression(ctx *IncDecExpressionContext) interface{}

	// Visit a parse tree produced by huloParser#preIncDecExpression.
	VisitPreIncDecExpression(ctx *PreIncDecExpressionContext) interface{}

	// Visit a parse tree produced by huloParser#postIncDecExpression.
	VisitPostIncDecExpression(ctx *PostIncDecExpressionContext) interface{}

	// Visit a parse tree produced by huloParser#factor.
	VisitFactor(ctx *FactorContext) interface{}

	// Visit a parse tree produced by huloParser#tripleExpression.
	VisitTripleExpression(ctx *TripleExpressionContext) interface{}

	// Visit a parse tree produced by huloParser#listExpression.
	VisitListExpression(ctx *ListExpressionContext) interface{}

	// Visit a parse tree produced by huloParser#newDelExpression.
	VisitNewDelExpression(ctx *NewDelExpressionContext) interface{}

	// Visit a parse tree produced by huloParser#mapExpression.
	VisitMapExpression(ctx *MapExpressionContext) interface{}

	// Visit a parse tree produced by huloParser#pair.
	VisitPair(ctx *PairContext) interface{}

	// Visit a parse tree produced by huloParser#variableExpression.
	VisitVariableExpression(ctx *VariableExpressionContext) interface{}

	// Visit a parse tree produced by huloParser#methodExpression.
	VisitMethodExpression(ctx *MethodExpressionContext) interface{}

	// Visit a parse tree produced by huloParser#expression.
	VisitExpression(ctx *ExpressionContext) interface{}

	// Visit a parse tree produced by huloParser#expressionList.
	VisitExpressionList(ctx *ExpressionListContext) interface{}

	// Visit a parse tree produced by huloParser#memberAccess.
	VisitMemberAccess(ctx *MemberAccessContext) interface{}

	// Visit a parse tree produced by huloParser#fileExpression.
	VisitFileExpression(ctx *FileExpressionContext) interface{}

	// Visit a parse tree produced by huloParser#callExpression.
	VisitCallExpression(ctx *CallExpressionContext) interface{}

	// Visit a parse tree produced by huloParser#callExpressionLinkedList.
	VisitCallExpressionLinkedList(ctx *CallExpressionLinkedListContext) interface{}

	// Visit a parse tree produced by huloParser#memberAccessPoint.
	VisitMemberAccessPoint(ctx *MemberAccessPointContext) interface{}

	// Visit a parse tree produced by huloParser#literal.
	VisitLiteral(ctx *LiteralContext) interface{}

	// Visit a parse tree produced by huloParser#receiverArgumentList.
	VisitReceiverArgumentList(ctx *ReceiverArgumentListContext) interface{}

	// Visit a parse tree produced by huloParser#namedArgumentList.
	VisitNamedArgumentList(ctx *NamedArgumentListContext) interface{}

	// Visit a parse tree produced by huloParser#namedArgument.
	VisitNamedArgument(ctx *NamedArgumentContext) interface{}

	// Visit a parse tree produced by huloParser#variableName.
	VisitVariableName(ctx *VariableNameContext) interface{}

	// Visit a parse tree produced by huloParser#rangeExpression.
	VisitRangeExpression(ctx *RangeExpressionContext) interface{}

	// Visit a parse tree produced by huloParser#expressionStatement.
	VisitExpressionStatement(ctx *ExpressionStatementContext) interface{}

	// Visit a parse tree produced by huloParser#option.
	VisitOption(ctx *OptionContext) interface{}

	// Visit a parse tree produced by huloParser#shortOption.
	VisitShortOption(ctx *ShortOptionContext) interface{}

	// Visit a parse tree produced by huloParser#longOption.
	VisitLongOption(ctx *LongOptionContext) interface{}

	// Visit a parse tree produced by huloParser#commandExpression.
	VisitCommandExpression(ctx *CommandExpressionContext) interface{}

	// Visit a parse tree produced by huloParser#commandJoin.
	VisitCommandJoin(ctx *CommandJoinContext) interface{}

	// Visit a parse tree produced by huloParser#commandStream.
	VisitCommandStream(ctx *CommandStreamContext) interface{}

	// Visit a parse tree produced by huloParser#commandMemberAccess.
	VisitCommandMemberAccess(ctx *CommandMemberAccessContext) interface{}

	// Visit a parse tree produced by huloParser#commandAccessPoint.
	VisitCommandAccessPoint(ctx *CommandAccessPointContext) interface{}

	// Visit a parse tree produced by huloParser#receiverParameters.
	VisitReceiverParameters(ctx *ReceiverParametersContext) interface{}

	// Visit a parse tree produced by huloParser#receiverParameterList.
	VisitReceiverParameterList(ctx *ReceiverParameterListContext) interface{}

	// Visit a parse tree produced by huloParser#receiverParameter.
	VisitReceiverParameter(ctx *ReceiverParameterContext) interface{}

	// Visit a parse tree produced by huloParser#namedParameters.
	VisitNamedParameters(ctx *NamedParametersContext) interface{}

	// Visit a parse tree produced by huloParser#namedParameterList.
	VisitNamedParameterList(ctx *NamedParameterListContext) interface{}

	// Visit a parse tree produced by huloParser#namedParameter.
	VisitNamedParameter(ctx *NamedParameterContext) interface{}

	// Visit a parse tree produced by huloParser#returnStatement.
	VisitReturnStatement(ctx *ReturnStatementContext) interface{}

	// Visit a parse tree produced by huloParser#functionDeclaration.
	VisitFunctionDeclaration(ctx *FunctionDeclarationContext) interface{}

	// Visit a parse tree produced by huloParser#operatorIdentifier.
	VisitOperatorIdentifier(ctx *OperatorIdentifierContext) interface{}

	// Visit a parse tree produced by huloParser#standardFunctionDeclaration.
	VisitStandardFunctionDeclaration(ctx *StandardFunctionDeclarationContext) interface{}

	// Visit a parse tree produced by huloParser#functionReturnValue.
	VisitFunctionReturnValue(ctx *FunctionReturnValueContext) interface{}

	// Visit a parse tree produced by huloParser#lambdaFunctionDeclaration.
	VisitLambdaFunctionDeclaration(ctx *LambdaFunctionDeclarationContext) interface{}

	// Visit a parse tree produced by huloParser#lambdaExpression.
	VisitLambdaExpression(ctx *LambdaExpressionContext) interface{}

	// Visit a parse tree produced by huloParser#lambdaBody.
	VisitLambdaBody(ctx *LambdaBodyContext) interface{}

	// Visit a parse tree produced by huloParser#functionSignature.
	VisitFunctionSignature(ctx *FunctionSignatureContext) interface{}

	// Visit a parse tree produced by huloParser#functionModifier.
	VisitFunctionModifier(ctx *FunctionModifierContext) interface{}

	// Visit a parse tree produced by huloParser#macroStatement.
	VisitMacroStatement(ctx *MacroStatementContext) interface{}

	// Visit a parse tree produced by huloParser#classDeclaration.
	VisitClassDeclaration(ctx *ClassDeclarationContext) interface{}

	// Visit a parse tree produced by huloParser#classModifier.
	VisitClassModifier(ctx *ClassModifierContext) interface{}

	// Visit a parse tree produced by huloParser#classSuper.
	VisitClassSuper(ctx *ClassSuperContext) interface{}

	// Visit a parse tree produced by huloParser#classBody.
	VisitClassBody(ctx *ClassBodyContext) interface{}

	// Visit a parse tree produced by huloParser#classMember.
	VisitClassMember(ctx *ClassMemberContext) interface{}

	// Visit a parse tree produced by huloParser#classMemberModifier.
	VisitClassMemberModifier(ctx *ClassMemberModifierContext) interface{}

	// Visit a parse tree produced by huloParser#classMethod.
	VisitClassMethod(ctx *ClassMethodContext) interface{}

	// Visit a parse tree produced by huloParser#classMethodModifier.
	VisitClassMethodModifier(ctx *ClassMethodModifierContext) interface{}

	// Visit a parse tree produced by huloParser#classBuiltinMethod.
	VisitClassBuiltinMethod(ctx *ClassBuiltinMethodContext) interface{}

	// Visit a parse tree produced by huloParser#classBuiltinMethodModifier.
	VisitClassBuiltinMethodModifier(ctx *ClassBuiltinMethodModifierContext) interface{}

	// Visit a parse tree produced by huloParser#classBuiltinParameters.
	VisitClassBuiltinParameters(ctx *ClassBuiltinParametersContext) interface{}

	// Visit a parse tree produced by huloParser#classNamedParameters.
	VisitClassNamedParameters(ctx *ClassNamedParametersContext) interface{}

	// Visit a parse tree produced by huloParser#classNamedParameterList.
	VisitClassNamedParameterList(ctx *ClassNamedParameterListContext) interface{}

	// Visit a parse tree produced by huloParser#classNamedParameter.
	VisitClassNamedParameter(ctx *ClassNamedParameterContext) interface{}

	// Visit a parse tree produced by huloParser#classNamedParameterAccessPoint.
	VisitClassNamedParameterAccessPoint(ctx *ClassNamedParameterAccessPointContext) interface{}

	// Visit a parse tree produced by huloParser#classInitializeExpression.
	VisitClassInitializeExpression(ctx *ClassInitializeExpressionContext) interface{}

	// Visit a parse tree produced by huloParser#enumDeclaration.
	VisitEnumDeclaration(ctx *EnumDeclarationContext) interface{}

	// Visit a parse tree produced by huloParser#enumModifier.
	VisitEnumModifier(ctx *EnumModifierContext) interface{}

	// Visit a parse tree produced by huloParser#enumBodySimple.
	VisitEnumBodySimple(ctx *EnumBodySimpleContext) interface{}

	// Visit a parse tree produced by huloParser#enumBodyAssociated.
	VisitEnumBodyAssociated(ctx *EnumBodyAssociatedContext) interface{}

	// Visit a parse tree produced by huloParser#enumBodyADT.
	VisitEnumBodyADT(ctx *EnumBodyADTContext) interface{}

	// Visit a parse tree produced by huloParser#enumValue.
	VisitEnumValue(ctx *EnumValueContext) interface{}

	// Visit a parse tree produced by huloParser#enumAssociatedFields.
	VisitEnumAssociatedFields(ctx *EnumAssociatedFieldsContext) interface{}

	// Visit a parse tree produced by huloParser#enumAssociatedField.
	VisitEnumAssociatedField(ctx *EnumAssociatedFieldContext) interface{}

	// Visit a parse tree produced by huloParser#enumFieldModifier.
	VisitEnumFieldModifier(ctx *EnumFieldModifierContext) interface{}

	// Visit a parse tree produced by huloParser#enumField.
	VisitEnumField(ctx *EnumFieldContext) interface{}

	// Visit a parse tree produced by huloParser#enumAssociatedValues.
	VisitEnumAssociatedValues(ctx *EnumAssociatedValuesContext) interface{}

	// Visit a parse tree produced by huloParser#enumAssociatedValue.
	VisitEnumAssociatedValue(ctx *EnumAssociatedValueContext) interface{}

	// Visit a parse tree produced by huloParser#enumAssociatedConstructor.
	VisitEnumAssociatedConstructor(ctx *EnumAssociatedConstructorContext) interface{}

	// Visit a parse tree produced by huloParser#enumConstructor.
	VisitEnumConstructor(ctx *EnumConstructorContext) interface{}

	// Visit a parse tree produced by huloParser#enumConstructorName.
	VisitEnumConstructorName(ctx *EnumConstructorNameContext) interface{}

	// Visit a parse tree produced by huloParser#enumConstructorParameters.
	VisitEnumConstructorParameters(ctx *EnumConstructorParametersContext) interface{}

	// Visit a parse tree produced by huloParser#enumConstructorDirectInit.
	VisitEnumConstructorDirectInit(ctx *EnumConstructorDirectInitContext) interface{}

	// Visit a parse tree produced by huloParser#enumConstructorInit.
	VisitEnumConstructorInit(ctx *EnumConstructorInitContext) interface{}

	// Visit a parse tree produced by huloParser#enumAssociatedMethods.
	VisitEnumAssociatedMethods(ctx *EnumAssociatedMethodsContext) interface{}

	// Visit a parse tree produced by huloParser#enumVariant.
	VisitEnumVariant(ctx *EnumVariantContext) interface{}

	// Visit a parse tree produced by huloParser#enumMethods.
	VisitEnumMethods(ctx *EnumMethodsContext) interface{}

	// Visit a parse tree produced by huloParser#enumMethod.
	VisitEnumMethod(ctx *EnumMethodContext) interface{}

	// Visit a parse tree produced by huloParser#enumMethodModifier.
	VisitEnumMethodModifier(ctx *EnumMethodModifierContext) interface{}

	// Visit a parse tree produced by huloParser#enumMember.
	VisitEnumMember(ctx *EnumMemberContext) interface{}

	// Visit a parse tree produced by huloParser#enumMemberModifier.
	VisitEnumMemberModifier(ctx *EnumMemberModifierContext) interface{}

	// Visit a parse tree produced by huloParser#enumBuiltinMethodModifier.
	VisitEnumBuiltinMethodModifier(ctx *EnumBuiltinMethodModifierContext) interface{}

	// Visit a parse tree produced by huloParser#enumInitialize.
	VisitEnumInitialize(ctx *EnumInitializeContext) interface{}

	// Visit a parse tree produced by huloParser#enumInitializeMember.
	VisitEnumInitializeMember(ctx *EnumInitializeMemberContext) interface{}

	// Visit a parse tree produced by huloParser#traitDeclaration.
	VisitTraitDeclaration(ctx *TraitDeclarationContext) interface{}

	// Visit a parse tree produced by huloParser#traitModifier.
	VisitTraitModifier(ctx *TraitModifierContext) interface{}

	// Visit a parse tree produced by huloParser#traitBody.
	VisitTraitBody(ctx *TraitBodyContext) interface{}

	// Visit a parse tree produced by huloParser#traitMember.
	VisitTraitMember(ctx *TraitMemberContext) interface{}

	// Visit a parse tree produced by huloParser#traitMemberModifier.
	VisitTraitMemberModifier(ctx *TraitMemberModifierContext) interface{}

	// Visit a parse tree produced by huloParser#implDeclaration.
	VisitImplDeclaration(ctx *ImplDeclarationContext) interface{}

	// Visit a parse tree produced by huloParser#implDeclarationBinding.
	VisitImplDeclarationBinding(ctx *ImplDeclarationBindingContext) interface{}

	// Visit a parse tree produced by huloParser#implDeclarationBody.
	VisitImplDeclarationBody(ctx *ImplDeclarationBodyContext) interface{}

	// Visit a parse tree produced by huloParser#extendDeclaration.
	VisitExtendDeclaration(ctx *ExtendDeclarationContext) interface{}

	// Visit a parse tree produced by huloParser#extendEnum.
	VisitExtendEnum(ctx *ExtendEnumContext) interface{}

	// Visit a parse tree produced by huloParser#extendClass.
	VisitExtendClass(ctx *ExtendClassContext) interface{}

	// Visit a parse tree produced by huloParser#extendTrait.
	VisitExtendTrait(ctx *ExtendTraitContext) interface{}

	// Visit a parse tree produced by huloParser#extendType.
	VisitExtendType(ctx *ExtendTypeContext) interface{}

	// Visit a parse tree produced by huloParser#extendMod.
	VisitExtendMod(ctx *ExtendModContext) interface{}

	// Visit a parse tree produced by huloParser#importDeclaration.
	VisitImportDeclaration(ctx *ImportDeclarationContext) interface{}

	// Visit a parse tree produced by huloParser#importSingle.
	VisitImportSingle(ctx *ImportSingleContext) interface{}

	// Visit a parse tree produced by huloParser#importAll.
	VisitImportAll(ctx *ImportAllContext) interface{}

	// Visit a parse tree produced by huloParser#importMulti.
	VisitImportMulti(ctx *ImportMultiContext) interface{}

	// Visit a parse tree produced by huloParser#asIdentifier.
	VisitAsIdentifier(ctx *AsIdentifierContext) interface{}

	// Visit a parse tree produced by huloParser#identifierAsIdentifier.
	VisitIdentifierAsIdentifier(ctx *IdentifierAsIdentifierContext) interface{}

	// Visit a parse tree produced by huloParser#moduleDeclaration.
	VisitModuleDeclaration(ctx *ModuleDeclarationContext) interface{}

	// Visit a parse tree produced by huloParser#moduleStatement.
	VisitModuleStatement(ctx *ModuleStatementContext) interface{}

	// Visit a parse tree produced by huloParser#useDeclaration.
	VisitUseDeclaration(ctx *UseDeclarationContext) interface{}

	// Visit a parse tree produced by huloParser#useSingle.
	VisitUseSingle(ctx *UseSingleContext) interface{}

	// Visit a parse tree produced by huloParser#useMulti.
	VisitUseMulti(ctx *UseMultiContext) interface{}

	// Visit a parse tree produced by huloParser#useAll.
	VisitUseAll(ctx *UseAllContext) interface{}

	// Visit a parse tree produced by huloParser#typeDeclaration.
	VisitTypeDeclaration(ctx *TypeDeclarationContext) interface{}

	// Visit a parse tree produced by huloParser#genericArguments.
	VisitGenericArguments(ctx *GenericArgumentsContext) interface{}

	// Visit a parse tree produced by huloParser#genericParameters.
	VisitGenericParameters(ctx *GenericParametersContext) interface{}

	// Visit a parse tree produced by huloParser#genericParameterList.
	VisitGenericParameterList(ctx *GenericParameterListContext) interface{}

	// Visit a parse tree produced by huloParser#genericParameter.
	VisitGenericParameter(ctx *GenericParameterContext) interface{}

	// Visit a parse tree produced by huloParser#compositeType.
	VisitCompositeType(ctx *CompositeTypeContext) interface{}

	// Visit a parse tree produced by huloParser#type.
	VisitType(ctx *TypeContext) interface{}

	// Visit a parse tree produced by huloParser#typeLiteral.
	VisitTypeLiteral(ctx *TypeLiteralContext) interface{}

	// Visit a parse tree produced by huloParser#nullableType.
	VisitNullableType(ctx *NullableTypeContext) interface{}

	// Visit a parse tree produced by huloParser#unionType.
	VisitUnionType(ctx *UnionTypeContext) interface{}

	// Visit a parse tree produced by huloParser#intersectionType.
	VisitIntersectionType(ctx *IntersectionTypeContext) interface{}

	// Visit a parse tree produced by huloParser#arrayType.
	VisitArrayType(ctx *ArrayTypeContext) interface{}

	// Visit a parse tree produced by huloParser#typeAccessPoint.
	VisitTypeAccessPoint(ctx *TypeAccessPointContext) interface{}

	// Visit a parse tree produced by huloParser#typeList.
	VisitTypeList(ctx *TypeListContext) interface{}

	// Visit a parse tree produced by huloParser#typeofExpression.
	VisitTypeofExpression(ctx *TypeofExpressionContext) interface{}

	// Visit a parse tree produced by huloParser#asExpression.
	VisitAsExpression(ctx *AsExpressionContext) interface{}

	// Visit a parse tree produced by huloParser#objectType.
	VisitObjectType(ctx *ObjectTypeContext) interface{}

	// Visit a parse tree produced by huloParser#objectTypeMember.
	VisitObjectTypeMember(ctx *ObjectTypeMemberContext) interface{}

	// Visit a parse tree produced by huloParser#tupleType.
	VisitTupleType(ctx *TupleTypeContext) interface{}

	// Visit a parse tree produced by huloParser#functionType.
	VisitFunctionType(ctx *FunctionTypeContext) interface{}

	// Visit a parse tree produced by huloParser#tryStatement.
	VisitTryStatement(ctx *TryStatementContext) interface{}

	// Visit a parse tree produced by huloParser#catchClause.
	VisitCatchClause(ctx *CatchClauseContext) interface{}

	// Visit a parse tree produced by huloParser#catchClauseReceiver.
	VisitCatchClauseReceiver(ctx *CatchClauseReceiverContext) interface{}

	// Visit a parse tree produced by huloParser#finallyClause.
	VisitFinallyClause(ctx *FinallyClauseContext) interface{}

	// Visit a parse tree produced by huloParser#throwStatement.
	VisitThrowStatement(ctx *ThrowStatementContext) interface{}

	// Visit a parse tree produced by huloParser#breakStatement.
	VisitBreakStatement(ctx *BreakStatementContext) interface{}

	// Visit a parse tree produced by huloParser#continueStatement.
	VisitContinueStatement(ctx *ContinueStatementContext) interface{}

	// Visit a parse tree produced by huloParser#ifStatement.
	VisitIfStatement(ctx *IfStatementContext) interface{}

	// Visit a parse tree produced by huloParser#matchStatement.
	VisitMatchStatement(ctx *MatchStatementContext) interface{}

	// Visit a parse tree produced by huloParser#matchCaseClause.
	VisitMatchCaseClause(ctx *MatchCaseClauseContext) interface{}

	// Visit a parse tree produced by huloParser#matchEnum.
	VisitMatchEnum(ctx *MatchEnumContext) interface{}

	// Visit a parse tree produced by huloParser#matchEnumMember.
	VisitMatchEnumMember(ctx *MatchEnumMemberContext) interface{}

	// Visit a parse tree produced by huloParser#matchTriple.
	VisitMatchTriple(ctx *MatchTripleContext) interface{}

	// Visit a parse tree produced by huloParser#matchTripleValue.
	VisitMatchTripleValue(ctx *MatchTripleValueContext) interface{}

	// Visit a parse tree produced by huloParser#matchCaseBody.
	VisitMatchCaseBody(ctx *MatchCaseBodyContext) interface{}

	// Visit a parse tree produced by huloParser#matchDefaultClause.
	VisitMatchDefaultClause(ctx *MatchDefaultClauseContext) interface{}

	// Visit a parse tree produced by huloParser#loopStatement.
	VisitLoopStatement(ctx *LoopStatementContext) interface{}

	// Visit a parse tree produced by huloParser#loopLabel.
	VisitLoopLabel(ctx *LoopLabelContext) interface{}

	// Visit a parse tree produced by huloParser#foreachStatement.
	VisitForeachStatement(ctx *ForeachStatementContext) interface{}

	// Visit a parse tree produced by huloParser#foreachClause.
	VisitForeachClause(ctx *ForeachClauseContext) interface{}

	// Visit a parse tree produced by huloParser#forStatement.
	VisitForStatement(ctx *ForStatementContext) interface{}

	// Visit a parse tree produced by huloParser#forClause.
	VisitForClause(ctx *ForClauseContext) interface{}

	// Visit a parse tree produced by huloParser#rangeStatement.
	VisitRangeStatement(ctx *RangeStatementContext) interface{}

	// Visit a parse tree produced by huloParser#rangeClause.
	VisitRangeClause(ctx *RangeClauseContext) interface{}

	// Visit a parse tree produced by huloParser#doWhileStatement.
	VisitDoWhileStatement(ctx *DoWhileStatementContext) interface{}

	// Visit a parse tree produced by huloParser#whileStatement.
	VisitWhileStatement(ctx *WhileStatementContext) interface{}

	// Visit a parse tree produced by huloParser#deferStatement.
	VisitDeferStatement(ctx *DeferStatementContext) interface{}

	// Visit a parse tree produced by huloParser#declareStatement.
	VisitDeclareStatement(ctx *DeclareStatementContext) interface{}

	// Visit a parse tree produced by huloParser#channelInputStatement.
	VisitChannelInputStatement(ctx *ChannelInputStatementContext) interface{}

	// Visit a parse tree produced by huloParser#channelOutputExpression.
	VisitChannelOutputExpression(ctx *ChannelOutputExpressionContext) interface{}

	// Visit a parse tree produced by huloParser#unsafeExpression.
	VisitUnsafeExpression(ctx *UnsafeExpressionContext) interface{}

	// Visit a parse tree produced by huloParser#comptimeExpression.
	VisitComptimeExpression(ctx *ComptimeExpressionContext) interface{}

	// Visit a parse tree produced by huloParser#externDeclaration.
	VisitExternDeclaration(ctx *ExternDeclarationContext) interface{}

	// Visit a parse tree produced by huloParser#externList.
	VisitExternList(ctx *ExternListContext) interface{}

	// Visit a parse tree produced by huloParser#externItem.
	VisitExternItem(ctx *ExternItemContext) interface{}
}
