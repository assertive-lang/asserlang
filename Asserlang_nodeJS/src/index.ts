import ReadLine from "readline-sync"
import { program } from "commander"
import path from "path"
import fs from "fs"

let codesSrc: string[] = []
const variables: { [key: string]: any } = {}
const localVariables: { [key: string]: { [key: string]: any } } = {}
const subRoutines: { [key: string]: (args?: any[]) => void | any } = {}
const statements = [
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
]

const execute = async (code: string) => {
  const lines: string[] = code.replace(/\r/gi, "").split("\n")
  codesSrc = [...lines]
  if (lines.shift() !== "쿠쿠루삥뽕") throw new Error("아무것도 모르죠?")
  if (lines.pop() !== "슉슈슉슉") throw new Error("아무것도 모르죠?")
  run(lines)
}

const run = async (lines: string[]) => {
  for (let line in lines) {
    const components = getComponents(lines[line])
    if (components.doesStartWithKeyword) {
      switch (components.keyword) {
        case ";;":
          const targetLine = toNumber(components.values.join("").trim())
          if (codesSrc[targetLine - 1]) {
            run(codesSrc.slice(targetLine - 1, codesSrc.length - 1))
            // const jumpingLines = lines.slice(targetLine, Number(line))
            // console.log(jumpingLines)
            // console.log([
            //   ...lines.slice(0, targetLine),
            //   ...jumpingLines,
            //   ...jumpingLines,
            //   lines[line] + "ㅋㅋ",
            //   ...lines.slice(Number(line) + 1, lines.length)
            // ])
          } else {
            throw new Error("어쩔GOTO인덱스;;")
          }
          break
        case "어쩔":
          declareVariable(lines[line])
          break
        case "저쩔":
          assignVariable(lines[line])
          break
        case "ㅇㅉ":
          print(lines[line])
          break
        case "ㅌㅂ":
          input(lines[line])
          break
        case "안물":
          const endIndex = lines.findIndex((v: any, i: any) => i !== Number(line) && v === "안물")
          if (endIndex <= -1) throw new Error("안물")
          const functionBlock = lines.splice(Number(line), endIndex)
          declareFunction(functionBlock)
          break
        case "안궁":
          callFunction(lines[line])
          break
        case "화났쥬?":
          conditionOperator(lines[line])
          break
        case "우짤래미":
          declareString(lines[line])
          break
        case "저짤래미":
          assignString(lines[line])
          break
      }
    }
  }
}

const getComponents = (
  line: string
):
  | {
      doesStartWithKeyword: false
      value: any
    }
  | {
      doesStartWithKeyword: true
      keyword: string
      values: string[]
    } => {
  const statement = statements.find((v) => line.startsWith(v))
  if (!statement) {
    if (isPureNumber(line)) return { doesStartWithKeyword: false, value: toNumber(line) }
    else
      return {
        doesStartWithKeyword: false,
        value: getVariable(line)
      }
  } else {
    return {
      doesStartWithKeyword: true,
      keyword: statement,
      values: line.trim().replace(statement, "").split("~")
    }
  }
}

const isPureNumber = (value: string) => {
  const add = value.split("").filter((v) => v === "ㅋ").length
  const subtract = value.split("").filter((v) => v === "ㅎ").length
  const index = value.split("").filter((v) => v === "ㅌ").length
  if (add + subtract + index === value.length) return true
  else return false
}

const toNumber = (line: string) => {
  if (!line.includes("ㅌ")) {
    const pluses = line
      .trim()
      .split("")
      .filter((v) => v === "ㅋ").length
    const minuses = line
      .trim()
      .split("")
      .filter((v) => v === "ㅎ").length
    return pluses - minuses
  }
  return Number(
    eval(
      line
        .trim()
        .split("ㅌ")
        .map((number) => {
          const pluses = number.split("").filter((v) => v === "ㅋ").length
          const minuses = number.split("").filter((v) => v === "ㅎ").length
          return pluses - minuses
        })
        .filter((v) => v)
        .join("*")
        .trim()
    )
  )
}

const toUnicode = (line: string) => {
  return String.fromCharCode(toNumber(line))
}

const getVariable = (line: string) => {
  if (line.startsWith("ㅌㅂ")) return input(line)
  if (line.startsWith("안궁")) return callFunction(line) ?? null
  if (statements.includes(line)) return
  if (isPureNumber(line)) return toNumber(line).toString()
  if (variables[line]) return variables[line]
  return null
}

const callFunction = (line: string): any => {
  if (!line.startsWith("안궁")) return
  const [name, ...args] = line.replace("안궁", "").split("~")
  return subRoutines[name](args.map((v) => getVariable(v)))
}

const declareFunction = (functionLines: string[]) => {
  const [declare, ...lines] = functionLines
  const declareComponents = getComponents(declare.trim())
  if (!declareComponents.doesStartWithKeyword) return
  if (declareComponents.keyword !== "안물") return
  const [name, ...args] = declareComponents.values
  if (!name) throw new Error("안물함수이름")
  localVariables[name] = {}
  const functionComponents: string[] = [`([${args.map((arg) => `${arg},`).join(" ")}]) => {`]
  for (const line of lines) {
    const lineComponents = getComponents(line)
    if (!lineComponents.doesStartWithKeyword) return
    if (lineComponents.keyword === "무지개반사") {
      const returnVariable = lineComponents.values.join("~").trim()
      functionComponents.push(
        `return ${
          args.includes(returnVariable)
            ? returnVariable
            : `${localVariables[name][returnVariable]}` ?? returnVariable.startsWith("ㅌㅂ")
            ? `input("${returnVariable}")`
            : `"${getVariable(returnVariable)}"`
        }`
      )
      break
    } else if (lineComponents.keyword === "어쩔" || lineComponents.keyword === "저쩔") {
      const [varName, ...varValue] = lineComponents.values
      if (varValue[0] === "ㅌㅂ") {
        functionComponents.push(
          `localVariables.${name}.${varName} = input("${varValue.join(" ").trim()}")`
        )
      } else {
        const value = varValue.join("~").trim()
        functionComponents.push(
          `localVariables.${name}.${varName} = ${
            isPureNumber(value)
              ? toNumber(value)
              : localVariables[name][value] ?? variables[value] ?? 0
          }`
        )
        localVariables[name][varName] = `${
          isPureNumber(value)
            ? toNumber(value)
            : localVariables[name][value] ?? variables[value] ?? 0
        }`
      }
    } else if (lineComponents.keyword === "ㅇㅉ") {
      functionComponents.push(
        `console.log(String(${lineComponents.values
          .map((value) =>
            args.includes(value)
              ? value
              : `"${
                  !isPureNumber(value)
                    ? localVariables[name][value] ?? getVariable(value)
                    : toNumber(value)
                }"`
          )
          .join(",")}))`
      )
    } else if (lineComponents.keyword === "ㅌㅂ") {
      functionComponents.push(`input("${lineComponents.values.join(" ")}`)
    }
  }
  functionComponents.push("}")
  localVariables[name] = {}
  subRoutines[name] = eval(functionComponents.join("\n"))
}

const conditionOperator = (line: string) => {
  const components = getComponents(line)
  if (!components.doesStartWithKeyword) return
  if (components.keyword !== "화났쥬?") return
  if (components.values.join("").trim().indexOf("킹받쥬?") <= -1) throw new Error("어쩔조건")
  const [condition, ...codes] = components.values.join("~").trim().split("킹받쥬?")
  const conditionValue = getVariable(condition)
  if (String(conditionValue) === "0") {
    run([codes.join("~").trim()])
  } else {
    return
  }
}

const declareVariable = (line: string) => {
  const components = getComponents(line)
  if (!components.doesStartWithKeyword) return
  if (components.keyword !== "어쩔") return
  const [name, first, ...values] = components.values
  if (!name || name.length <= 0 || statements.includes(name) || isPureNumber(name))
    throw new Error("어쩔변수이름")
  let allocatingValue = ""
  if (first) {
    if (first === "ㅌㅂ") {
      const inputValue = input(first + values.join("~"))
      if (inputValue) allocatingValue = inputValue
    } else if (isPureNumber(first)) {
      allocatingValue = String(toNumber(first + values.join("~")))
    } else {
      allocatingValue = "0"
    }
  } else {
    allocatingValue = "0"
  }
  variables[name] = allocatingValue
}

const assignVariable = (line: string) => {
  const components = getComponents(line)
  if (!components.doesStartWithKeyword) return
  if (components.keyword !== "저쩔") return
  const [name, ...values] = components.values
  if (!name || name.length <= 0) throw new Error("어쩔변수")
  const doesVariableExist = getVariable(name)
  if (doesVariableExist === null) throw new Error("어쩔변수")
  let value = ""
  if (name === "ㅌㅂ") {
    const inputValue = input("ㅌㅂ" + values.join("~").trim())
    if (inputValue) value = inputValue
  } else {
    value = String(toNumber(values.join("~").trim()))
  }
  variables[name] = value
}

const declareString = (line: string) => {
  const components = getComponents(line)
  if (!components.doesStartWithKeyword) return
  if (components.keyword !== "우짤래미") return
  const [name, first, ...values] = components.values
  if (!name || name.length <= 0 || statements.includes(name) || isPureNumber(name))
    throw new Error("어쩔변수이름")
  let allocatingValue = ""
  if (first) {
    if (first === "ㅌㅂ") {
      const inputValue = input(first + values.join("~"))
      if (inputValue) allocatingValue = inputValue
    } else if (isPureNumber(first)) {
      allocatingValue = String(toUnicode(first + values.join("~")))
    } else {
      allocatingValue = String(toUnicode("0"))
    }
  } else {
    allocatingValue = "0"
  }
  variables[name] = allocatingValue
}

const assignString = (line: string) => {
  const components = getComponents(line)
  if (!components.doesStartWithKeyword) return
  if (components.keyword !== "저짤래미") return
  const [name, ...values] = components.values
  if (!name || name.length <= 0) throw new Error("어쩔변수")
  const doesVariableExist = getVariable(name)
  if (doesVariableExist === null) throw new Error("어쩔변수")
  let value = ""
  if (name === "ㅌㅂ") {
    const inputValue = input("ㅌㅂ" + values.join("~").trim())
    if (inputValue) value = inputValue
  } else {
    value = String(toUnicode(values.join("~").trim()))
  }
  variables[name] = value
}

const print = (line: string) => {
  const components = getComponents(line)
  if (!components.doesStartWithKeyword) return
  if (components.keyword !== "ㅇㅉ") return
  console.log(
    line
      .trim()
      .replace("ㅇㅉ", "")
      .split(" ")
      .map((v) => getVariable(v))
      .join("~")
  )
}

const input = (line: string) => {
  const components = getComponents(line)
  if (!components.doesStartWithKeyword) return
  if (components.keyword !== "ㅌㅂ") return
  const inputUser = ReadLine.question(components.values.join(" ") + "\n", { encoding: "utf-8" })
  return inputUser
}

program.parse()

const targetFilePath = path.join(process.cwd(), program.args.join(" "))
if (!fs.existsSync(targetFilePath)) {
  throw new Error("어쩔파일")
}
const codes = fs.readFileSync(targetFilePath)
execute(codes.toString("utf-8"))
