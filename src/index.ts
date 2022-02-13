import ReadLine from "readline-sync"
import { program } from "commander"
import path from "path"
import fs from "fs"

const variables: { [key: string]: any } = {}
const localVariables: { [key: string]: { [key: string]: any } } = {}
const subRoutines: { [key: string]: (args?: any[]) => void | any } = {}
const statements = ["ㅇㅉ", "ㅌㅂ", "저쩔", "어쩔", "안물", "안물", "안궁"]

const execute = async (code: string) => {
  const lines: string[] = code.replace(/\r/gi, "").split("\n")
  if (lines.shift() !== "쿠쿠루삥뽕") throw new Error("아무것도 모르죠?")
  if (lines.pop() !== "슉슈슉슉") throw new Error("아무것도 모르죠?")
  run(lines)
}

const run = async (lines: any) => {
  for (const line in lines) {
    const [statement] = lines[line].split(" ")
    if (statement === "어쩔") {
      declareVariable(lines[line])
    } else if (statement === "저쩔") {
      assignVariable(lines[line])
    } else if (statement === "ㅇㅉ") {
      print(lines[line])
    } else if (statement === "ㅌㅂ") {
      input(lines[line])
    } else if (statement === "안물") {
      const endIndex = lines.findIndex((v: any, i: any) => i !== Number(line) && v === "안물")
      if (endIndex <= -1) throw new Error("안물")
      const functionBlock = lines.splice(Number(line), endIndex)
      declareFunction(functionBlock)
    } else if (statement === "안궁") {
      callFunction(lines[line])
    } else if (statement === "화났쥬?") {
      conditionOperator(lines[line])
    }
  }
}

const isPureNumber = (value: string) => {
  const ㅋ = value.split("").filter((v) => v === "ㅋ").length
  const ㅎ = value.split("").filter((v) => v === "ㅎ").length
  if (ㅋ + ㅎ === value.length) return true
  else return false
}

const getVariable = (line: string) => {
  if (line.split(" ")[0] === "ㅌㅂ") return input(line)
  if (statements.includes(line.split(" ")[0])) return
  const value = variables[line]
  if (!value) return toNumber(line) === 0 ? null : toNumber(line)
  return value
}

const callFunction = (line: string) => {
  const components = line.split(" ")
  if (components.shift() !== "안궁") return
  const [name, ...args] = components
  subRoutines[name](args.map((v) => getVariable(v)))
}

const declareFunction = (functionLines: string[]) => {
  const [declare, ...lines] = functionLines
  const [_, name, ...args] = declare.split(" ")
  if (!name) throw new Error("안물함수이름")
  localVariables[name] = {}
  const functionComponents: string[] = [`([${args.map((arg) => `${arg},`).join(" ")}]) => {`]
  for (const line of lines) {
    const [statement, ...values] = line.trim().split(" ")
    if (statement === "어쩔" || statement === "저쩔") {
      const [varName, ...varValue] = values
      if (varValue[0] === "ㅌㅂ") {
        functionComponents.push(
          `localVariables.${name}.${varName} = input("${varValue.join(" ").trim()}")`
        )
      } else {
        const value = varValue.join(" ").trim()
        functionComponents.push(
          `localVariables.${name}.${varName} = ${
            isPureNumber(value) // todo: number to string
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
    } else if (statement === "ㅇㅉ") {
      functionComponents.push(
        `console.log(${values
          .map((value) =>
            args.includes(value)
              ? value
              : `"${
                  !isPureNumber(value)
                    ? localVariables[name][value] ?? getVariable(value)
                    : toNumber(value)
                }"`
          )
          .join(",")})`
      )
    } else if (statement === "ㅌㅂ") {
      functionComponents.push(`input("${values.join(" ")}`)
    }
  }
  functionComponents.push("}")
  localVariables[name] = {}
  subRoutines[name] = eval(functionComponents.join("\n"))
}

const conditionOperator = (line: string) => {
  let [statement, condition, isTrue, ...values]: string[] = line.split(" ")
  if (statement !== "화났쥬?") return
  const conditionValue = getVariable(condition)
  if (conditionValue === "0") {
    if (isTrue === "킹받쥬?") {
      run(values.join(" "))
    } else {
      throw new Error("어쩔조건")
    }
  } else {
    return
  }
}

const declareVariable = (line: string) => {
  const [statement, name, first, ...values] = line.split(" ")
  if (statement !== "어쩔") return
  if (!name || name.length <= 0 || statements.includes(name) || isPureNumber(name))
    throw new Error("어쩔변수")
  let allocatingValue = ""
  if (first) {
    if (first === "ㅌㅂ") {
      const inputValue = input([first, ...values].join(" "))
      if (inputValue) allocatingValue = inputValue
    } else if (isPureNumber(first)) {
      allocatingValue = String(toNumber([first, ...values].join(" ")))
    } else {
      allocatingValue = "0"
    }
  } else {
    allocatingValue = "0"
  }
  variables[name] = allocatingValue
}

const assignVariable = (line: string) => {
  const [statement, name, ...values] = line.split(" ")
  if (statement !== "저쩔") return
  if (!name || name.length <= 0) throw new Error("어쩔변수")
  const doesVariableExist = getVariable(name)
  if (doesVariableExist === null) throw new Error("어쩔변수")
  let value = ""
  if (name === "ㅌㅂ") {
    const inputValue = input(values.join(" "))
    if (inputValue) value = inputValue
  } else {
    value = String(toNumber(values.join(" ")))
  }
  variables[name] = value
}

const toNumber = (line: string) => {
  const pluses = line.split("").filter((v) => v === "ㅋ").length
  const minuses = line.split("").filter((v) => v === "ㅎ").length
  return pluses - minuses
}

const print = (line: string) => {
  const [statement, ...printOut] = line.split(" ")
  if (statement !== "ㅇㅉ") return
  console.log(printOut.map((v) => getVariable(v)).join(" "))
}

const input = (line: string) => {
  const [statement, ...inputString] = line.split(" ")
  if (statement !== "ㅌㅂ") return
  const inputUser = ReadLine.question(inputString.join(" ") + "\n")
  return inputUser
}

program.parse()

const targetFilePath = path.join(process.cwd(), program.args.join(" "))
if (!fs.existsSync(targetFilePath)) {
  throw new Error("어쩔파일")
}
const codes = fs.readFileSync(targetFilePath)
execute(codes.toString("utf-8"))
