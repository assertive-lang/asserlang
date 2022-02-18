using System;
using System.IO;
using System.Linq;
using System.Reflection;
using System.Reflection.Emit;
using System.Collections.Generic;

namespace AssertiveLang
{
    public static class Program
    {
        public static void Main(string[] args)
        {
            new AsserLangCompiler(File.ReadAllText(args[0])).Compile().AsserFunction();
        }
    }
    public class AsserLangCompiler
    {
        public AsserLangCompiler(string source)
        {
            var arr = source.Split(new[] { "\r\n", "\r", "\n" }, StringSplitOptions.None);
            if (arr.First() != "쿠쿠루삥뽕") throw new AsserLangException($"아무것도 모르죠?");
            if (arr.Last() != "슉슈슉슉") throw new AsserLangException($"아무것도 모르죠?");
            Source = source;
            TypeBuilder = AsserLangMethod.AsserLangModule.DefineType($"AsserType{(AsserTypeCount++ == 1 ? "" : AsserTypeCount.ToString())}", TypeAttributes.Public);
            Method = (AsserLangMethod)TypeBuilder.DefineMethod("AsserLangMethod", MethodAttributes.Public | MethodAttributes.Static, typeof(object), Type.EmptyTypes);
            AsserLangMethod.AsserLangAsm.SetEntryPoint(Method.Method);
        }
        public string Source { get; }
        public TypeBuilder TypeBuilder { get; }
        public AsserLangMethod Method { get; private set; }
        public ILGenerator IL { get; }
        public Dictionary<string, LocalBuilder> Variables = new Dictionary<string, LocalBuilder>();
        public Dictionary<string, AsserLangMethod> Methods = new Dictionary<string, AsserLangMethod>();
        public List<int> GotoList = new List<int>();
        public AsserLangResult Result;
        public AsserLangResult Compile()
        {
            int lineNumber = 0;
            using (StringReader reader = new StringReader(Source))
            {
                string line = "";
                while ((line = reader.ReadLine()) != null)
                {
                    lineNumber++;
                    if (line.TrimStart().StartsWith("//"))
                        continue;
                    if (line.Contains(";;"))
                    {
                        line = line.Substring(line.LastIndexOf(";;")).Replace(";;", "");
                        if (AsserLangMethod.IsNumber(line, out int label))
                            GotoList.Add(label);
                    }
                }
                reader.Close();
            }
            Method.PrepareLabels(GotoList);
            lineNumber = 0;
            try
            {
                using (StringReader reader = new StringReader(Source))
                {
                    string line = "";
                    while ((line = reader.ReadLine()) != null)
                    {
                        lineNumber++;
                        if (line.TrimStart().StartsWith("//"))
                            continue;
                        Method.Emit(line, reader, TypeBuilder, ref lineNumber);
                    }
                }
                if (!Method.Returned)
                    Method.Ret(true);
                var created = TypeBuilder.CreateType();
                var method = created.GetMethod(Method.Method.Name);
                var func = (Func<object>)method.CreateDelegate(typeof(Func<object>));
                return Result = new AsserLangResult(created, method, func);
            }
            catch (Exception ex)
            {
                throw new AsserLangException($"컴파일 오류! 줄 {lineNumber}", ex);
            }
        }
        public static void Save()
            => AsserLangMethod.AsserLangAsm.Save("AsserLangMethod.dll");
        public static int AsserTypeCount = 1;
    }
    public class AsserLangException : Exception
    {
        public AsserLangException(string msg) : base(msg) { }
        public AsserLangException(string msg, Exception inner) : base(msg, inner) { }
    }
    public class AsserLangResult
    {
        public Type AsserType;
        public MethodInfo AsserLangMethod;
        public Func<object> AsserFunction;
        public AsserLangResult(Type asserType, MethodInfo asserMethod, Func<object> asserFunction)
        {
            AsserType = asserType;
            AsserFunction = asserFunction;
            AsserLangMethod = asserMethod;
        }
    }
    public class AsserLangMethod
    {
        public static explicit operator AsserLangMethod(MethodBuilder methodBuilder) => new AsserLangMethod(methodBuilder);
        public static explicit operator MethodBuilder(AsserLangMethod asserMethod) => asserMethod.Method;
        public static MethodInfo sEquality = typeof(string).GetMethod("op_Equality");
        public static MethodInfo iParse = typeof(int).GetMethod("Parse", new[] { typeof(string) });
        public static MethodInfo oWrite = typeof(Console).GetMethod("Write", new[] { typeof(object) });
        public static MethodInfo ReadL = typeof(Console).GetMethod("ReadLine");
        public MethodBuilder Method;
        public ILGenerator IL;
        public Dictionary<string, short> Parameters;
        public Dictionary<string, LocalBuilder> Variables;
        public Dictionary<LocalBuilder, Type> BoxType;
        public Dictionary<short, Type> BoxTypeArg;
        public AsserLangMethod Outer;
        public Type LastBoxedType = typeof(int);
        public string[] ParameterNames;
        public bool RequireLdnull = true;
        public bool Returned = false;
        public Stack<Label> Conditions;
        public Dictionary<int, Label> Labels = new Dictionary<int, Label>();
        public Dictionary<string, AsserLangMethod> LocalMethods;
        public AsserLangMethod(MethodBuilder method, Dictionary<string, short> parameters = null, string[] parameterNames = null, Dictionary<string, AsserLangMethod> localMethods = null, AsserLangMethod outer = null)
        {
            Method = method;
            IL = method.GetILGenerator();
            Parameters = parameters ?? new Dictionary<string, short>();
            ParameterNames = parameterNames;
            Variables = new Dictionary<string, LocalBuilder>();
            BoxType = new Dictionary<LocalBuilder, Type>();
            BoxTypeArg = new Dictionary<short, Type>();
            LocalMethods = localMethods ?? new Dictionary<string, AsserLangMethod>();
            Labels = new Dictionary<int, Label>();
            Conditions = new Stack<Label>();
            Outer = outer;
        }
        public AsserLangMethod DefineMethod(TypeBuilder typeBuilder, string name, string[] parameterNames)
        {
            var meth = typeBuilder.DefineMethod(name, MethodAttributes.Public | MethodAttributes.Static, typeof(object), Enumerable.Repeat(typeof(object), parameterNames.Length).ToArray());
            Dictionary<string, short> param = new Dictionary<string, short>();
            for (int i = 0; i < parameterNames.Length; i++)
            {
                param.Add(parameterNames[i], (short)i);
                meth.DefineParameter(i + 1, ParameterAttributes.None, parameterNames[i]);
            }
            var locMeth = new AsserLangMethod(meth, param, parameterNames, null, this);
            locMeth.PrepareLabels(Labels.Keys.ToList());
            LocalMethods.Add(name, locMeth);
            return locMeth;
        }
        public void Box(Type type)
        {
            Emit(OpCodes.Box, type);
            LastBoxedType = type;
        }
        public void PrepareLabels(List<int> gotoList)
        {
            gotoList.ForEach(i =>
            {
                if (!Labels.ContainsKey(i))
                    Labels.Add(i, IL.DefineLabel());
            });
        }
        public void StartCondition(string name, bool equals = true)
        {
            RequireLdnull = true;
            Get(name);
            if (LastBoxedType == typeof(int))
            {
                Emit(OpCodes.Ldc_I4, 0);
                Emit(OpCodes.Ceq);
            }
            else if (LastBoxedType == typeof(string))
            {
                Emit(OpCodes.Ldstr, "0");
                Emit(OpCodes.Call, sEquality);
            }
            var label = IL.DefineLabel();
            Emit(equals ? OpCodes.Brfalse : OpCodes.Brtrue, label);
            Conditions.Push(label);
        }
        public void EndCondition()
        {
            if (!Conditions.Any()) return;
            RequireLdnull = true;
            IL.MarkLabel(Conditions.Pop());
        }
        public void Goto(Label label)
        {
            RequireLdnull = true;
            Emit(OpCodes.Br, label);
        }
        public LocalBuilder DeclareLocal(string name, object value = null)
        {
            RequireLdnull = true;
            var loc = IL.DeclareLocal(typeof(object));
            if (name != null)
            {
                loc.SetLocalSymInfo(name);
                Variables.Add(name, loc);
            }
            if (value is int i)
            {
                Emit(OpCodes.Ldc_I4, i);
                Emit(OpCodes.Box, typeof(int));
                BoxType[loc] = typeof(int);
            }
            else if (value is string s)
            {
                Emit(OpCodes.Ldstr, s);
                Emit(OpCodes.Box, typeof(string));
                BoxType[loc] = typeof(string);
            }
            else
            {
                Emit(OpCodes.Ldc_I4_0);
                Emit(OpCodes.Box, typeof(int));
                BoxType[loc] = typeof(int);
            }
            Emit(OpCodes.Stloc, loc);
            LastBoxedType = BoxType[loc];
            return loc;
        }
        public void Get(string name, bool ignoreBoxing = false)
        {
            if (Variables.TryGetValue(name, out _))
                GetLocal(name, ignoreBoxing);
            else if (Parameters.TryGetValue(name, out _))
                GetParam(name, ignoreBoxing);
            else throw new AsserLangException($"어쩔변수 ({name}이(가) 선언되어있지 않음)");
        }
        public void Set(string name, object value = null)
        {
            if (Variables.TryGetValue(name, out _))
                SetLocal(name, value);
            else if (Parameters.TryGetValue(name, out _))
                SetParam(name, value);
            else throw new AsserLangException($"어쩔변수 ({name}이(가) 선언되어있지 않음)");
        }
        public void GetLocal(string name, bool ignoreBoxing = false)
        {
            RequireLdnull = false;
            Emit(OpCodes.Ldloc, Variables[name]);
            Type boxType = BoxType[Variables[name]];
            if (!ignoreBoxing)
                Emit(OpCodes.Unbox_Any, boxType);
            LastBoxedType = boxType;
        }
        public void SetLocal(string name, object value)
        {
            RequireLdnull = true;
            var loc = Variables[name];
            if (value is int i)
            {
                Emit(OpCodes.Ldc_I4, i);
                Emit(OpCodes.Box, typeof(int));
                BoxType[loc] = typeof(int);
            }
            else if (value is string s)
            {
                Emit(OpCodes.Ldstr, s);
                Emit(OpCodes.Box, typeof(string));
                BoxType[loc] = typeof(string);
            }
            Emit(OpCodes.Stloc, loc);
            LastBoxedType = BoxType[loc];
        }
        public void GetParam(string name, bool ignoreBoxing = false)
        {
            RequireLdnull = false;
            Emit(OpCodes.Ldarg, Parameters[name]);
            if (BoxTypeArg.TryGetValue(Parameters[name], out Type boxType))
            {
                if (!ignoreBoxing)
                    Emit(OpCodes.Unbox_Any, boxType);
                LastBoxedType = boxType;
            }
            else
            {
                if (!ignoreBoxing)
                    Emit(OpCodes.Unbox_Any, typeof(int));
                BoxTypeArg[Parameters[name]] = typeof(int);
                LastBoxedType = typeof(int);
            }
        }
        public void SetParam(string name, object value)
        {
            RequireLdnull = true;
            var arg = Parameters[name];
            if (value is int i)
            {
                Emit(OpCodes.Ldc_I4, i);
                Emit(OpCodes.Box, typeof(int));
                BoxTypeArg[arg] = typeof(int);
            }
            else if (value is string s)
            {
                Emit(OpCodes.Ldstr, s);
                Emit(OpCodes.Box, typeof(string));
                BoxTypeArg[arg] = typeof(string);
            }
            else throw new AsserLangException($"Compilation Error! Invalid LocalVariable Type ({value.GetType()})");
            Emit(OpCodes.Starg, arg);
            LastBoxedType = BoxTypeArg[arg];
        }
        public Label DefineAndMark()
        {
            RequireLdnull = true;
            var label = IL.DefineLabel();
            IL.MarkLabel(label);
            return label;
        }
        public void Ret(bool returnLocal = false)
        {
            if (!Returned)
            {
                if (RequireLdnull)
                    IL.Emit(OpCodes.Ldnull);
                else IL.Emit(OpCodes.Box, LastBoxedType);
                IL.Emit(OpCodes.Ret);
            }
            Returned = true;
            if (returnLocal)
                LocalMethods.Values.ToList().ForEach(m => m.Ret());
        }
        public void Emit(string line, TextReader reader, TypeBuilder typeBuilder, ref int lineNumber)
        {
            if (Labels.ContainsKey(lineNumber))
                IL.MarkLabel(Labels[lineNumber]);
            switch (GetKeyword(line, out var values))
            {
                case "어쩔":
                    if (values.Length > 1)
                    {
                        string k;
                        if ((k = GetKeyword(values[1], out string[] valss)) != null)
                        {
                            DeclareLocal(values[0]);
                            Emit(InstructionToReEmit(k, valss), reader, typeBuilder, ref lineNumber);
                            Set(values[0]);
                        }
                        else if (IsNumber(values[1], out int result))
                            DeclareLocal(values[0], result);
                        else if (ChkVariableOperation(values[1]))
                        {
                            DeclareLocal(values[0], result);
                            ComputeVariableOperation(values[0], values[1]);
                        }
                        else
                        {
                            DeclareLocal(values[0], result);
                            Set(values[0], values[1]);
                        }
                    }
                    else DeclareLocal(values[0]);
                    break;
                case "저쩔":
                    string kk;
                    if ((kk = GetKeyword(values[1], out string[] vals)) != null)
                    {
                        Emit(InstructionToReEmit(kk, vals), reader, typeBuilder, ref lineNumber);
                        Set(values[0]);
                    }
                    else if (IsNumber(values[1], out int re))
                        Set(values[0], re);
                    else if (ChkVariableOperation(values[1]))
                        ComputeVariableOperation(values[0], values[1]);
                    else Set(values[0], values[1]);
                    break;
                case "ㅇㅉ":
                    if ((kk = GetKeyword(values[0], out vals)) != null)
                    {
                        LastBoxedType = typeof(string);
                        Emit(InstructionToReEmit(kk, vals), reader, typeBuilder, ref lineNumber);
                    }
                    else
                    {
                        if (IsNumber(values[0], out int result))
                        {
                            Emit(OpCodes.Ldc_I4, result);
                            Box(typeof(int));
                        }
                        else if (ChkVariableOperation(values[0]))
                            ComputeVariableOperation(null, values[0], true);
                        else
                            Get(values[0], true);
                    }
                    Emit(OpCodes.Call, oWrite);
                    RequireLdnull = true;
                    break;
                case "ㅌㅂ":
                    Emit(OpCodes.Call, ReadL);
                    Emit(OpCodes.Box, typeof(string));
                    RequireLdnull = false;
                    break;
                case "안물":
                    var name = values[0];
                    if (string.IsNullOrWhiteSpace(name))
                        throw new AsserLangException("안물 (잘못된 함수 이름)");
                    string[] parameters = RemoveStart(values, 1);
                    var localNew = DefineMethod(typeBuilder, name, parameters);
                    while (!(line = reader.ReadLine()).StartsWith("안물"))
                    {
                        lineNumber++;
                        localNew.Emit(line, reader, typeBuilder, ref lineNumber);
                    }
                    lineNumber++;
                    break;
                case "안궁":
                    if (LocalMethods.TryGetValue(values[0], out var toCall))
                    {
                        for (int i = 1; i < values.Length; i++)
                        {
                            var param = values[i];
                            if (GetKeyword(param, out _) != null)
                                Emit(param, reader, typeBuilder, ref lineNumber);
                            else if (IsNumber(param, out int result))
                            {
                                Emit(OpCodes.Ldc_I4, result);
                                Box(typeof(int));
                            }
                            else Get(values[i]);
                        }
                        RequireLdnull = false;
                        Emit(OpCodes.Call, toCall.Method);
                    }
                    else throw new AsserLangException($"안궁 (\"{values[0]}\"이름의 함수는 선언되지 않았습니다.)");
                    break;
                case "화났쥬?":
                    var sp = values[0].Split('?');
                    var inst = sp[0].Replace("킹받쥬", "");
                    StartCondition(inst);
                    string kw;
                    if ((kw = GetKeyword(sp[1], out vals)) != null)
                        Emit(InstructionToReEmit(kw, vals), reader, typeBuilder, ref lineNumber);
                    else throw new AsserLangException("어쩔조건 (조건식이 존재하지 않습니다.)");
                    EndCondition();
                    break;
                case "킹받쥬?":
                    sp = values[0].Split('?');
                    inst = sp[0].Replace("화났쥬", "");
                    StartCondition(inst, false);
                    if ((kk = GetKeyword(sp[1], out vals)) != null)
                        Emit(InstructionToReEmit(kk, vals), reader, typeBuilder, ref lineNumber);
                    else throw new AsserLangException("어쩔조건 (조건식이 존재하지 않습니다.)");
                    EndCondition();
                    break;
                case "우짤래미":
                    var decTo = 0;
                    if (values.Length > 1)
                    {
                        if (IsNumber(values[1], out var result))
                            decTo = result;
                        else throw new AsserLangException($"어쩔변수 ({values[1]}은(는) 숫자가 아닙니다)");
                    }
                    DeclareLocal(values[0], char.ConvertFromUtf32(decTo));
                    break;
                case "저짤래미":
                    if (IsNumber(values[1], out int res))
                        Set(values[0], char.ConvertFromUtf32(res));
                    else throw new AsserLangException($"어쩔변수 ({values[1]}은(는) 숫자가 아닙니다)");
                    break;
                case "무지개반사":
                    if (values.Length > 0)
                    {
                        Get(values[0]);
                    }
                    Ret();
                    break;
                case ";;":
                    if (IsNumber(values[0], out int lineNum))
                        Goto(Labels[lineNum]);
                    else throw new AsserLangException($"어쩔GOTO인덱스;; (인덱스 \"{values[0]}\"은(는) 숫자가 아닙니다)");
                    break;
            }
        }
        private void Emit(OpCode opcode, params object[] operands)
        {
            if (operands.Length == 1)
                switch (operands[0])
                {
                    case string i:
                        IL.Emit(opcode, i);
                        return;
                    case FieldInfo i:
                        IL.Emit(opcode, i);
                        return;
                    case Label[] i:
                        IL.Emit(opcode, i);
                        return;
                    case Label i:
                        IL.Emit(opcode, i);
                        return;
                    case LocalBuilder i:
                        IL.Emit(opcode, i);
                        return;
                    case float i:
                        IL.Emit(opcode, i);
                        return;
                    case byte i:
                        IL.Emit(opcode, i);
                        return;
                    case sbyte i:
                        IL.Emit(opcode, i);
                        return;
                    case short i:
                        IL.Emit(opcode, i);
                        return;
                    case double i:
                        IL.Emit(opcode, i);
                        return;
                    case MethodInfo i:
                        RequireLdnull = i.ReturnType == typeof(void);
                        IL.Emit(opcode, i);
                        return;
                    case int i:
                        IL.Emit(opcode, i);
                        return;
                    case long i:
                        IL.Emit(opcode, i);
                        return;
                    case Type i:
                        IL.Emit(opcode, i);
                        return;
                    case SignatureHelper i:
                        IL.Emit(opcode, i);
                        return;
                    case ConstructorInfo i:
                        IL.Emit(opcode, i);
                        return;
                    default:
                        IL.Emit(opcode);
                        return;
                }
            else if (operands.Length == 2)
                switch (operands[0])
                {
                    case MethodInfo i:
                        RequireLdnull = i.ReturnType == typeof(void);
                        IL.EmitCall(opcode, i, (Type[])operands[1]);
                        return;
                    default:
                        throw new InvalidOperationException();
                }
            else if (operands.Length == 3)
                switch (operands[0])
                {
                    case System.Runtime.InteropServices.CallingConvention i:
                        IL.EmitCalli(opcode, i, (Type)operands[1], (Type[])operands[2]);
                        return;
                    default:
                        throw new InvalidOperationException();
                }
            else if (operands.Length == 4)
                switch (operands[0])
                {
                    case CallingConventions i:
                        IL.EmitCalli(opcode, i, (Type)operands[1], (Type[])operands[2], (Type[])operands[3]);
                        return;
                    default:
                        throw new InvalidOperationException();
                }
            else
                IL.Emit(opcode);
        }
        private void ComputeVariableOperation(string name, string value, bool box = false)
        {
            var split = value.Split('ㅌ');
            if (split.Length > 1)
            {
                split.Aggregate((cur, next) =>
                {
                    int curP = cur.Where(c => c == 'ㅋ').Count();
                    int curM = cur.Where(c => c == 'ㅎ').Count();

                    int nextP = next.Where(c => c == 'ㅋ').Count();
                    int nextM = next.Where(c => c == 'ㅎ').Count();

                    var curVar = cur.Replace("ㅋ", "").Replace("ㅎ", "");
                    var nextVar = next.Replace("ㅋ", "").Replace("ㅎ", "");

                    if (curP - curM != 0)
                    {
                        Get(curVar);
                        if (LastBoxedType == typeof(string))
                        {
                            Emit(OpCodes.Unbox, LastBoxedType);
                            Emit(OpCodes.Call, iParse);
                        }
                        Emit(OpCodes.Ldc_I4, curP - curM);
                        Emit(OpCodes.Add);
                    }
                    else Get(curVar);

                    if (nextP - nextM != 0)
                    {
                        Get(nextVar);
                        if (LastBoxedType == typeof(string))
                        {
                            Emit(OpCodes.Unbox, LastBoxedType);
                            Emit(OpCodes.Call, iParse);
                        }
                        Emit(OpCodes.Ldc_I4, nextP - nextM);
                        Emit(OpCodes.Add);
                    }
                    else Get(nextVar);

                    Emit(OpCodes.Mul);
                    if (box)
                        Box(typeof(int));
                    if (!string.IsNullOrWhiteSpace(name))
                    {
                        if (!box)
                            Box(typeof(int));
                        Set(name);
                    }
                    return null;
                });
            }
            else
            {
                int p = value.Where(c => c == 'ㅋ').Count();
                int m = value.Where(c => c == 'ㅎ').Count();
                var var = value.Replace("ㅋ", "").Replace("ㅎ", "");
                if (p - m != 0)
                {
                    Get(var);
                    if (LastBoxedType == typeof(string))
                    {
                        Emit(OpCodes.Unbox, LastBoxedType);
                        Emit(OpCodes.Call, iParse);
                    }
                    Emit(OpCodes.Ldc_I4, p - m);
                    Emit(OpCodes.Add);
                }
                else Get(var);
                if (box)
                    Box(typeof(int));
                if (!string.IsNullOrWhiteSpace(name))
                {
                    if (!box)
                        Box(typeof(int));
                    Set(name);
                }
            }
        }
        #region Statics
        static AsserLangMethod()
        {
            if (Type.GetType("Mono.Runtime") != null)
            {
                AsserLangAsm = AssemblyBuilder.DefineDynamicAssembly(new AssemblyName("AsserLangMethod"), AssemblyBuilderAccess.Run);
                AsserLangModule = AsserLangAsm.DefineDynamicModule("AsserLangMethod");
            }
            else
            {
                AsserLangAsm = AssemblyBuilder.DefineDynamicAssembly(new AssemblyName("AsserLangMethod"), AssemblyBuilderAccess.RunAndSave);
                AsserLangModule = AsserLangAsm.DefineDynamicModule("AsserLangMethod", "AsserLangMethod.dll", true);
            }
        }
        public static AssemblyBuilder AsserLangAsm;
        public static ModuleBuilder AsserLangModule;
        public static bool IsNumber(string value, out int result)
        {
            if (!string.IsNullOrWhiteSpace(value.Replace("ㅋ", "").Replace("ㅎ", "").Replace("ㅌ", "")))
            {
                result = 0;
                return false;
            }
            return int.TryParse(value.Trim().Split('ㅌ').Select(s =>
            {
                int pluses = s.Where(c => c == 'ㅋ').Count();
                int minuses = s.Where(c => c == 'ㅎ').Count();
                return pluses - minuses;
            }).Aggregate((cur, next) => cur * next).ToString(), out result);
        }
        public static bool ChkVariableOperation(string value)
            => !string.IsNullOrWhiteSpace(value.Replace("ㅋ", "").Replace("ㅎ", "").Replace("ㅌ", "")) && (value.Contains("ㅋ") || value.Contains("ㅎ") || value.Contains("ㅌ") && !value.Contains("ㅌㅂ"));
        public static string GetKeyword(string line, out string[] values)
        {
            for (int i = 0; i < States.Length; i++)
            {
                var state = States[i];
                if (line.StartsWith(state))
                {
                    if (state == "화났쥬?" || state == "킹받쥬?" || state == "ㅇㅉ")
                    {
                        values = new string[1];
                        values[0] = line.Replace(state, "");
                    }
                    else values = line.Replace(state, "").Split('~');
                    return state;
                }
            }
            values = null;
            return null;
        }
        public static T[] RemoveStart<T>(T[] array, int count)
        {
            T[] arr = new T[array.Length - count];
            Array.Copy(array, count, arr, 0, arr.Length);
            return arr;
        }
        public static T[] RemoveEnd<T>(T[] array, int count)
        {
            T[] arr = new T[array.Length - count];
            Array.Copy(array, 0, arr, 0, arr.Length);
            return arr;
        }
        public static string InstructionToReEmit(string keyword, string[] vals)
        {
            string seed = $"{keyword}{vals[0]}";
            return RemoveStart(vals, 1).Aggregate(seed, (cur, next) => $"{cur}~{next}");
        }
        public static readonly string[] States = new string[13]
        {
            "ㅇㅉ",
            "ㅌㅂ",
            "저쩔",
            "어쩔",
            "안물",
            "안물",
            "안궁",
            "화났쥬?",
            "킹받쥬?",
            "우짤래미",
            "저짤래미",
            "무지개반사",
            ";;"
        };
        #endregion
    }
}
