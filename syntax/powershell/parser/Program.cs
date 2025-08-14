using System.Management.Automation;
using System.Management.Automation.Language;

namespace PowerShellToHulo
{
    class Program
    {
        static void Main(string[] args)
        {
            if (args.Length < 1)
            {
                Console.WriteLine("Usage: PowerShellToHulo <input.ps1> [output.hl]");
                return;
            }

            string inputFile = args[0];
            string outputFile = args.Length > 1 ? args[1] : Path.ChangeExtension(inputFile, ".hl");

            try
            {
                string psScript = File.ReadAllText(inputFile);
                string huloCode = ConvertPowerShellToHulo(psScript);
                File.WriteAllText(outputFile, huloCode);
                Console.WriteLine($"Converted {inputFile} to {outputFile}");
            }
            catch (Exception ex)
            {
                Console.WriteLine($"Error: {ex.Message}");
            }
        }

        static string ConvertPowerShellToHulo(string psScript)
        {
            // 解析 PowerShell 脚本
            Token[] tokens;
            ParseError[] errors;
            Ast ast = Parser.ParseInput(psScript, out tokens, out errors);

            if (errors.Length > 0)
            {
                throw new Exception($"PowerShell parse errors: {string.Join(", ", errors.Select(e => e.Message))}");
            }

            // 转换为 Hulo 代码
            var converter = new PowerShellToHuloConverter();
            return converter.Convert(ast);
        }
    }

    class PowerShellToHuloConverter
    {
        public string Convert(Ast ast)
        {
            var huloCode = new System.Text.StringBuilder();

            // 遍历 AST 并转换
            var allNodes = ast.FindAll(x => true, true);

            foreach (var node in allNodes)
            {
                switch (node)
                {
                    case FunctionDefinitionAst func:
                        ConvertFunction(func, huloCode);
                        break;
                    case AssignmentStatementAst assignment:
                        ConvertAssignment(assignment, huloCode);
                        break;
                    case IfStatementAst ifStmt:
                        ConvertIfStatement(ifStmt, huloCode);
                        break;
                    case ForStatementAst forStmt:
                        ConvertForStatement(forStmt, huloCode);
                        break;
                    case WhileStatementAst whileStmt:
                        ConvertWhileStatement(whileStmt, huloCode);
                        break;
                    case CommandAst cmd:
                        ConvertCommand(cmd, huloCode);
                        break;
                }
            }

            return huloCode.ToString();
        }

        private void ConvertFunction(FunctionDefinitionAst func, System.Text.StringBuilder sb)
        {
            sb.AppendLine($"func {func.Name}(");

            if (func.Parameters != null)
            {
                var paramsList = func.Parameters.Select(p => p.Name);
                sb.AppendLine(string.Join(", ", paramsList));
            }

            sb.AppendLine(") {");

            if (func.Body != null)
            {
                // 递归转换函数体
                Convert(func.Body);
            }

            sb.AppendLine("}");
        }

        private void ConvertAssignment(AssignmentStatementAst assignment, System.Text.StringBuilder sb)
        {
            string left = ConvertExpression(assignment.Left);
            string right = ConvertExpression(assignment.Right);
            sb.AppendLine($"  {left} = {right}");
        }

                private void ConvertIfStatement(IfStatementAst ifStmt, System.Text.StringBuilder sb)
        {
            // PowerShell IfStatementAst 的结构不同
            if (ifStmt.Clauses != null && ifStmt.Clauses.Count > 0)
            {
                var firstClause = ifStmt.Clauses[0];
                string condition = ConvertExpression(firstClause.Item1);
                sb.AppendLine($"  if {condition} {{");

                if (firstClause.Item2 != null)
                {
                    Convert(firstClause.Item2);
                }

                // 处理 else 子句
                if (ifStmt.Clauses.Count > 1)
                {
                    sb.AppendLine("  } else {");
                    var elseClause = ifStmt.Clauses[1];
                    if (elseClause.Item2 != null)
                    {
                        Convert(elseClause.Item2);
                    }
                }

                sb.AppendLine("  }");
            }
        }

        private void ConvertForStatement(ForStatementAst forStmt, System.Text.StringBuilder sb)
        {
            sb.AppendLine("  for {");
            // 转换 for 循环
            sb.AppendLine("  }");
        }

        private void ConvertWhileStatement(WhileStatementAst whileStmt, System.Text.StringBuilder sb)
        {
            string condition = ConvertExpression(whileStmt.Condition);
            sb.AppendLine($"  while {condition} {{");

            if (whileStmt.Body != null)
            {
                Convert(whileStmt.Body);
            }

            sb.AppendLine("  }");
        }

        private void ConvertCommand(CommandAst cmd, System.Text.StringBuilder sb)
        {
            string commandName = cmd.CommandElements.FirstOrDefault()?.ToString() ?? "";
            var arguments = cmd.CommandElements.Skip(1).Select(ConvertExpression);

            sb.AppendLine($"  {commandName}({string.Join(", ", arguments)})");
        }

        private string ConvertExpression(Ast expression)
        {
            // 这里需要根据具体的表达式类型进行转换
            return expression?.ToString() ?? "";
        }
    }
}
