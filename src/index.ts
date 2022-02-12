import ReadLine from "readline-sync"
import { program } from "commander"
import path from "path"
import fs from "fs"

const variables: { [key: string]: any } = {}
const localVariables: { [key: string]: any } = {}
const statements = ["ㅇㅉ", "ㅌㅂ", "저쩔", "어쩔", "어쩔함수", "저쩔함수"]

const execute = async (code: string) => {
  const lines: string[] = code.replace(/\r/gi, "").split("\n")
  if (lines.shift() !== "쿠쿠루삥뽕") throw new Error("아무것도 모르죠?")
  if (lines.pop() !== "슉슈슉슉") throw new Error("아무것도 모르죠?")
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
    }/* else if (statement === "어쩔함수") {
      const endIndex = lines.slice(Number(line)).indexOf("저쩔함수")
      if (endIndex <= -1) throw new Error("저쩔함수")
      const functionBlock = lines.slice(Number(line), endIndex - 1)
      parseFunction(functionBlock)
    }*/
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

const parseFunction = (functionLines: string[]) => {
  console.log(functionLines)
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
