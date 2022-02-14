using System;
using System.IO;
using System.Linq;
using System.Reflection;
using System.Reflection.Emit;
using System.Collections.Generic;

namespace AssertiveLang
{
    public class Program
    {
        public static void Main(string[] args)
        {
            AsserLang.Load(args[0]).Execute();
        }
    }
    public class AsserLang
    {
        static AsserLang()
        {
            if (Type.GetType("Mono.Runtime") != null)
            {
                AsserLangAsm = AssemblyBuilder.DefineDynamicAssembly(new AssemblyName("AsserLang"), AssemblyBuilderAccess.Run);
                AsserLangModule = AsserLangAsm.DefineDynamicModule("AsserLang");
            }
            else
            {
                AsserLangAsm = AssemblyBuilder.DefineDynamicAssembly(new AssemblyName("AsserLang"), AssemblyBuilderAccess.RunAndSave);
                AsserLangModule = AsserLangAsm.DefineDynamicModule("AsserLang", "AsserLang.dll", true);
            }
        }
        public static AssemblyBuilder AsserLangAsm;
        public static ModuleBuilder AsserLangModule;
        public static MethodInfo Console_WriteLine = typeof(Console).GetMethod("WriteLine", new[] { typeof(string) });
        public static MethodInfo Console_ReadLine = typeof(Console).GetMethod("ReadLine");
        public static MethodInfo String_Equality = typeof(string).GetMethod("op_Equality");
        public static int AsserCount = 0;
        public static readonly string[] States = new string[12]
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
            "무지개반사"
        };
        public string Source;
        public MethodInfo AsserMethod;
        public Action AsserFunction;
        public Type AsserType;
        public AsserLang(string source)
        {
            Source = source;
            var arr = source.Split(new[] { "\r\n", "\r", "\n" }, StringSplitOptions.None);
            if (arr.First() != "쿠쿠루삥뽕") throw new AsserLangException("아무것도 모르죠?");
            if (arr.Last() != "슉슈슉슉") throw new AsserLangException("아무것도 모르죠?");
        }
        public void Compile()
        {
            TypeBuilder asserType = AsserLangModule.DefineType($"AsserType{(AsserCount++ == 0 ? "" : AsserCount.ToString())}", TypeAttributes.Public);
            MethodBuilder asserMethod = asserType.DefineMethod("AsserMethod", MethodAttributes.Public | MethodAttributes.Static, CallingConventions.Standard, typeof(void), Type.EmptyTypes);
            var ilGen = asserMethod.GetILGenerator();
            Dictionary<string, short> localVars = new Dictionary<string, short>();
            Dictionary<string, Dictionary<string, short>> localMethVars = new Dictionary<string, Dictionary<string, short>>();
            Dictionary<string, Dictionary<string, short>> localMethArgs = new Dictionary<string, Dictionary<string, short>>();
            Dictionary<string, MethodBuilder> localMeths = new Dictionary<string, MethodBuilder>();
            bool cond = false;
            bool reqNull = false;
            using (StringReader sr = new StringReader(Source))
            {
                string line = "";
                int lineNumber = 1;
                while ((line = sr.ReadLine()) != null)
                {
                    try
                    {
                        if (GetKeyword(line, out var vals) == "안물")
                        {
                            try
                            {
                                MethodBuilder asserLocal = asserType.DefineMethod(vals[0], MethodAttributes.Public | MethodAttributes.Static, typeof(string), GetParameters(vals.Length - 1));
                                localMeths[vals[0]] = asserLocal;
                                localMethVars[vals[0]] = new Dictionary<string, short>();
                                localMethArgs[vals[0]] = new Dictionary<string, short>();
                                bool reqNullLoc = false;
                                for (int i = 0; i < vals.Length - 1; i++)
                                {
                                    localMethArgs[vals[0]].Add(vals[i + 1], (short)i);
                                    asserLocal.DefineParameter(i + 1, ParameterAttributes.None, vals[i + 1]);
                                }
                                var ilG = asserLocal.GetILGenerator();
                                while (!(line = sr.ReadLine()).StartsWith("안물"))
                                    Emit(ilG, line, lineNumber++, localMethVars[vals[0]], null, localMethArgs[vals[0]], ref cond, ref reqNullLoc);
                                Emit(ilG, line, lineNumber++, localVars, localMeths, null, ref cond, ref reqNullLoc);
                            }
                            catch (Exception ex)
                            {
                                throw new AsserLangException("안물", ex);
                            }
                        }
                        else Emit(ilGen, line, lineNumber++, localVars, localMeths, null, ref cond, ref reqNull);
                    }
                    catch (Exception ex)
                    {
                        throw new AsserLangException($"컴파일 오류! 줄 {lineNumber} ({ex.Message})", ex);
                    }
                }
                if (!reqNull)
                    ilGen.Emit(OpCodes.Pop);
                ilGen.Emit(OpCodes.Ret);
            }
            AsserType = asserType.CreateType();
            AsserMethod = AsserType.GetMethod("AsserMethod");
            AsserFunction = (Action)AsserMethod.CreateDelegate(typeof(Action));
        }
        public void Execute()
            => AsserFunction();
        public static AsserLang Load(string path)
        {
            if (!File.Exists(path))
                throw new AsserLangException("어쩔파일", new FileNotFoundException(path));
            AsserLang asserLang = new AsserLang(File.ReadAllText(path));
            asserLang.Compile();
            return asserLang;
        }
        public static void Emit(ILGenerator il, string line, int lineNumber, Dictionary<string, short> locVars, Dictionary<string, MethodBuilder> locMeths, Dictionary<string, short> locMethArgs, ref bool condition, ref bool requireLdnull)
        {
            var retLabel = default(Label?);
            switch (GetKeyword(line, out string[] values))
            {
                case "어쩔":
                    if (condition)
                    {
                        retLabel = il.DefineLabel();
                        il.Emit(OpCodes.Brfalse, (Label)retLabel);
                    }
                    var loc = il.DeclareLocal(typeof(string));
                    locVars[values[0]] = (short)loc.LocalIndex;
                    if (values.Length > 1)
                    {
                        string value = values[1];
                        if (IsNumber(value, out int number))
                            il.Emit(OpCodes.Ldstr, number.ToString());
                        else il.Emit(OpCodes.Ldstr, value);
                        il.Emit(OpCodes.Stloc, loc);
                    }
                    else
                    {
                        il.Emit(OpCodes.Ldstr, "0");
                        il.Emit(OpCodes.Stloc, loc);
                    }
                    if (condition)
                    {
                        il.MarkLabel((Label)retLabel);
                        condition = false;
                    }
                    requireLdnull = true;
                    break;
                case "저쩔":
                    if (condition)
                    {
                        retLabel = il.DefineLabel();
                        il.Emit(OpCodes.Brfalse, (Label)retLabel);
                    }
                    if (locVars.TryGetValue(values[0], out var locIdx))
                    {
                        string value = values[1];
                        if (IsNumber(value, out int number))
                            il.Emit(OpCodes.Ldstr, number.ToString());
                        else il.Emit(OpCodes.Ldstr, value);
                        il.Emit(OpCodes.Stloc, locIdx);
                    }
                    else if (locMethArgs.TryGetValue(values[0], out var argIdx))
                    {
                        string value = values[1];
                        if (IsNumber(value, out int number))
                            il.Emit(OpCodes.Ldstr, number.ToString());
                        else il.Emit(OpCodes.Ldstr, value);
                        il.Emit(OpCodes.Starg, argIdx);
                    }
                    else throw new AsserLangException("어쩔변수");
                    if (condition)
                    {
                        il.MarkLabel((Label)retLabel);
                        condition = false;
                    }
                    requireLdnull = true;
                    break;
                case "ㅇㅉ":
                    if (condition)
                    {
                        retLabel = il.DefineLabel();
                        il.Emit(OpCodes.Brfalse, (Label)retLabel);
                    }
                    if (GetKeyword(values[0], out string[] vals) != null)
                    {
                        Emit(il, values[0], lineNumber++, locVars, locMeths, locMethArgs, ref condition, ref requireLdnull);
                        il.Emit(OpCodes.Call, Console_WriteLine);
                    }
                    else
                    {
                        if (locVars.TryGetValue(values[0], out locIdx))
                            il.Emit(OpCodes.Ldloc, locIdx);
                        else if (locMethArgs.TryGetValue(values[0], out var argIdx))
                            il.Emit(OpCodes.Ldarg, argIdx);
                        else throw new AsserLangException("어쩔변수");
                        il.Emit(OpCodes.Call, Console_WriteLine);
                    }
                    if (condition)
                    {
                        il.MarkLabel((Label)retLabel);
                        condition = false;
                    }
                    requireLdnull = true;
                    break;
                case "ㅌㅂ":
                    if (condition)
                    {
                        retLabel = il.DefineLabel();
                        il.Emit(OpCodes.Brfalse, (Label)retLabel);
                    }
                    if (values.Length < 2)
                    {
                        il.Emit(OpCodes.Call, Console_ReadLine);
                        requireLdnull = false;
                    }
                    else
                    {
                        il.Emit(OpCodes.Call, Console_ReadLine);
                        if (locVars.TryGetValue(values[0], out locIdx))
                            il.Emit(OpCodes.Stloc, locIdx);
                        else if (locMethArgs.TryGetValue(values[0], out var argIdx))
                            il.Emit(OpCodes.Starg, argIdx);
                        else throw new AsserLangException("어쩔변수");
                        requireLdnull = true;
                    }
                    if (condition)
                    {
                        il.MarkLabel((Label)retLabel);
                        condition = false;
                    }
                    break;
                case "안물":
                    if (requireLdnull)
                        il.Emit(OpCodes.Ldnull);
                    il.Emit(OpCodes.Ret);
                    break;
                case "안궁":
                    if (condition)
                    {
                        retLabel = il.DefineLabel();
                        il.Emit(OpCodes.Brfalse, (Label)retLabel);
                    }
                    if (locMeths.TryGetValue(values[0], out var meth))
                    {
                        for (int i = 0; i < values.Length - 1; i++)
                        {
                            string val = values[i + 1];
                            if (IsNumber(val, out int number))
                                il.Emit(OpCodes.Ldstr, number.ToString());
                            else il.Emit(OpCodes.Ldstr, val);
                        }
                        il.Emit(OpCodes.Call, meth);
                    }
                    if (condition)
                    {
                        il.MarkLabel((Label)retLabel);
                        condition = false;
                    }
                    requireLdnull = false;
                    break;
                case "화났쥬?":
                    if (condition)
                    {
                        retLabel = il.DefineLabel();
                        il.Emit(OpCodes.Brfalse, (Label)retLabel);
                    }
                    string[] split = values[0].Split('?');
                    var n = split[0].Replace("킹받쥬", "");
                    if (locVars.TryGetValue(n, out locIdx))
                        il.Emit(OpCodes.Ldloc, locIdx);
                    else if (locMethArgs.TryGetValue(n, out var argIdx))
                        il.Emit(OpCodes.Ldarg, argIdx);
                    else throw new AsserLangException("어쩔조건");
                    il.Emit(OpCodes.Ldstr, "0");
                    il.Emit(OpCodes.Call, String_Equality);
                    condition = true;
                    var agg = values.Aggregate("", (cur, next) => $"{cur}~{next}");
                    var lastIdx = agg.LastIndexOf('?') + 1;
                    Emit(il, agg.Substring(lastIdx), lineNumber++, locVars, locMeths, locMethArgs, ref condition, ref requireLdnull);
                    if (condition)
                    {
                        il.MarkLabel((Label)retLabel);
                        condition = false;
                    }
                    requireLdnull = true;
                    break;
                case "우짤래미":
                    if (condition)
                    {
                        retLabel = il.DefineLabel();
                        il.Emit(OpCodes.Brfalse, (Label)retLabel);
                    }
                    loc = il.DeclareLocal(typeof(string));
                    locVars[values[0]] = (short)loc.LocalIndex;
                    if (values.Length > 1)
                    {
                        string value = values[1];
                        if (IsNumber(value, out int number))
                            il.Emit(OpCodes.Ldstr, char.ConvertFromUtf32((int)number));
                        else il.Emit(OpCodes.Ldstr, value);
                        il.Emit(OpCodes.Stloc, loc);
                    }
                    else
                    {
                        il.Emit(OpCodes.Ldstr, "0");
                        il.Emit(OpCodes.Stloc, loc);
                    }
                    if (condition)
                    {
                        il.MarkLabel((Label)retLabel);
                        condition = false;
                    }
                    requireLdnull = true;
                    break;
                case "저짤래미":
                    if (condition)
                    {
                        retLabel = il.DefineLabel();
                        il.Emit(OpCodes.Brfalse, (Label)retLabel);
                    }
                    if (locVars.TryGetValue(values[0], out locIdx))
                    {
                        string value = values[1];
                        if (IsNumber(value, out int number))
                            il.Emit(OpCodes.Ldstr, char.ConvertFromUtf32((int)number));
                        else il.Emit(OpCodes.Ldstr, value);
                        il.Emit(OpCodes.Stloc, locIdx);
                    }
                    else if (locMethArgs.TryGetValue(values[0], out var argIdx))
                    {
                        string value = values[1];
                        if (IsNumber(value, out int number))
                            il.Emit(OpCodes.Ldstr, char.ConvertFromUtf32((int)number));
                        else il.Emit(OpCodes.Ldstr, value);
                        il.Emit(OpCodes.Starg, argIdx);
                    }
                    else throw new AsserLangException("어쩔변수");
                    if (condition)
                    {
                        il.MarkLabel((Label)retLabel);
                        condition = false;
                    }
                    requireLdnull = true;
                    break;
                case "무지개반사":
                    if (condition)
                    {
                        retLabel = il.DefineLabel();
                        il.Emit(OpCodes.Brfalse, (Label)retLabel);
                    }
                    if (locVars.TryGetValue(values[0], out locIdx))
                        il.Emit(OpCodes.Ldloc, locIdx);
                    else if (locMethArgs.TryGetValue(values[0], out var argIdx))
                        il.Emit(OpCodes.Ldarg, argIdx);
                    else throw new AsserLangException("어쩔변수");
                    if (condition)
                    {
                        il.MarkLabel((Label)retLabel);
                        condition = false;
                    }
                    il.Emit(OpCodes.Ret);
                    requireLdnull = false;
                    break;
            }
        }
        public static void SaveAsserLangAsm()
            => AsserLangAsm.Save("AsserLang.dll");
        public static Type[] GetParameters(int length)
            => Enumerable.Repeat(typeof(string), length).ToArray();
        public static bool IsNumber(string value, out int result)
        {
            return int.TryParse(value.Trim().Split('ㅌ').Select(s =>
            {
                int pluses = s.Where(c => c == 'ㅋ').Count();
                int minuses = s.Where(c => c == 'ㅎ').Count();
                return pluses - minuses;
            }).Aggregate((cur, next) => cur * next).ToString(), out result);
        }
        public static string GetKeyword(string line, out string[] values)
        {
            for (int i = 0; i < States.Length; i++)
            {
                var state = States[i];
                if (line.StartsWith(state))
                {
                    values = line.Replace(state, "").Split('~');
                    return state;
                }
            }
            values = null;
            return null;
        }
    }
    public class AsserLangException : Exception
    {
        public AsserLangException(string msg) : base(msg) { }
        public AsserLangException(string msg, Exception inner) : base(msg, inner) { }
    }
}
