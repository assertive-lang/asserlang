import os
import sys
from typing import Dict, Optional, Union, List


def end_letter(word: str, yes: Optional[str] = "은", no: Optional[str] = "는") -> str:
    last_word = word[-1]
    if ord("ㄱ") <= ord(last_word) <= ord("ㅎ"):
        return yes
    if ord("ㅏ") <= ord(last_word) <= ord("ㅣ"):
        return no
    if last_word.lower() in "aeiouwy":
        return no
    elif last_word.isalpha():
        return yes
    if (ord(last_word) - ord("가")) % 28:
        return yes
    else:
        return no


class func:
    def __init__(self, name: str, start: int, param: Dict[str, int]):
        self.name = name
        self.start = start
        self.cnt = start
        self.var = param
        self.var_uni = {}


class asserlang:
    def __init__(self) -> None:
        self.version = "v1.5.3"
        self.keywords = ("ㅋ", "ㅎ", "ㅌ",
                         "어쩔", "저쩔", "우짤래미", "저짤래미",
                         "ㅇㅉ", "ㅌㅂ", "안물", "안궁",
                         "무지개반사", "화났쥬?", "킹받쥬?", ";;")
        self.exec = (None, None, None,
                     self.make_var, self.assign_var, self.make_var_uni, self.assign_var_uni,
                     self.print, None, self.make_func, self.call_func,
                     self.retn, self.condition, None, self.jump)
        self.call = {}
        self.funcs = [func("__main__", 1, {})]
        self.lines = []
        self.writing_func = False
        self.stop = False
        self.file = ""
        self.return_value = False

    def error(self, state: str) -> None:
        print("Traceback (most recent call last):")
        for i in self.funcs:
            print("    ", end="")
            if self.file:
                print(f"file \"{self.file}\", ", end="")
            print(f"line {i.cnt}, in {i.name}")
            print(f"      {self.lines[i.cnt-1]}")
        print(state)
        self.stop = True

    def calc(self, value: str) -> Union[None, int, str]:
        if value == "":
            return 0
        value = value.replace("ㅌㅂ", "ㅇㅉ")
        func_count = value.count("안궁")
        if func_count > 1:
            self.error("안물안궁: 한 줄에서 2번 이상 함수를 호출함")
            return None
        result_last = 0
        if func_count == 1:
            if self.return_value is False:
                line = value.split("안궁")[1]
                self.funcs[-1].cnt -= 1
                self.call_func(line)
                return None
            if self.return_value is None:
                return None
            result_last = self.return_value
            self.return_value = False
            value = value.split("안궁")[0]
        var = self.funcs[-1].var
        var.update({"ㅋ": 1, "ㅎ": -1})
        uni = self.funcs[-1].var_uni
        names = list(sorted(list(var.keys()) + list(uni.keys()),
                     key=lambda x: len(x), reverse=True))
        result = []
        digits = value.split("ㅌ")
        return_uni = False
        for index, i in enumerate(digits):
            result.append(0)
            while i:
                for j in names:
                    if i.startswith(j):
                        result[index] += var[j] if j in var else uni[j]
                        if j in uni:
                            return_uni = True
                        break
                else:
                    if i.startswith("ㅇㅉ"):
                        in_value = input("입력: ")
                        if not (in_value.isnumeric() or (in_value[0] == "-" and in_value[1:].isnumeric())):
                            self.error("어쩔ㅌㅂ: 입력으로 정수가 아닌 것이 입력됨")
                            return None
                        result[index] += int(in_value)
                        j = "ㅇㅉ"
                    else:
                        self.error("어쩔변수: 해당하는 변수가 없음")
                        return None
                i = i[len(j):]
        result[-1] += result_last
        result = eval('*'.join(str(n) for n in result))
        return chr(result) if return_uni else result

    def jump(self, line: str) -> None:
        value = self.calc(line)
        if value is None:
            return
        if type(value) is str:
            value = ord(value)
        if value == self.funcs[-1].cnt:
            self.error("어쩔;;;;: 자신이 위치한 줄로는 점프할 수 없음")
            return
        if self.funcs[-1].name == "__main__":
            if value > len(self.lines)+1:
                self.error("어쩔;;;;: 아직 작성되지 않은 부분으로 점프할 수 없음")
                return
            for i, j in self.call.items():
                if j[0] <= value <= j[2]:
                    self.error(f"어쩔;;;;: 함수 밖에서 함수 \"{i}\"{end_letter(i, '으로', '로')} 점프할 수 없음")
                    return
        else:
            j = self.call[self.funcs[-1].name]
            if value <= j[0] or j[1] <= value:
                self.error(f"어쩔;;;;: 함수 \"{self.funcs[-1].name}\" 안에서 함수 밖으로 점프할 수 없음")
                return
        self.funcs[-1].cnt = value-1

    def condition(self, line: str) -> None:
        cnt = line.count("킹받쥬?")
        if cnt == 0:
            self.error("어쩔조건: 킹받쥬?가 없음")
            return
        if cnt > 1:
            self.error(f"어쩔조건: 킹받쥬?가 {cnt}개 있음")
            return
        condition, line = line.split("킹받쥬?")
        value = self.calc(condition)
        if value is None:
            return
        if type(value) is str:
            value = ord(value)
        if value == 0:
            self.execute_line(line)

    def check_name(self, name: str) -> str:
        for i in self.keywords:
            if i in name:
                return i
        return ""

    def retn(self, line: Union[None, str]) -> None:
        if line:
            self.return_value = self.calc(line)
        elif line == "":
            self.return_value = 0
        elif line is None:
            self.return_value = None
        self.funcs.pop()

    def make_func(self, line: str) -> None:
        if line.strip("~") == "" and not self.writing_func:
            if len(self.funcs) == 1:
                self.error("안물안궁: 함수 이름이 필요함")
                return
            self.retn("")
            return
        if line == "":
            self.call[self.writing_func].append(self.funcs[-1].cnt)
            self.writing_func = False
            return
        if self.writing_func:
            self.error("안물안궁: 함수 안에서 함수를 작성함")
            del self.call[self.writing_func]
            self.writing_func = False
            return
        line = line.split("~")
        include = self.check_name(line[0])
        if include:
            self.error(f"안물안궁: 함수 \"{line[0]}\"{end_letter(line[0])} 키워드 \"{include}\"{end_letter(include, '을', '를')} 포함함")
            return
        names = []
        for name in line[1:]:
            include = self.check_name(name)
            if include:
                self.error(f"안물안궁: 매개변수 \"{name}\"{end_letter(name)} 키워드 \"{include}\"{end_letter(include, '을', '를')} 포함함")
                return
            if name in names:
                self.error(f"안물안궁: 매개변수 \"{name}\"{end_letter(name)} 다른 매개변수와 겹침")
                return
            names.append(name)
        self.call[line[0]] = [self.funcs[-1].cnt, names]
        self.writing_func = line[0]

    def call_func(self, line: str) -> None:
        if line.strip("~") == "":
            self.error("안물안궁: 함수 이름이 필요함")
            return None
        line = line.split("~")
        if line[0] not in self.call:
            self.error(f"안물안궁: \"{line[0]}\"{end_letter(line[0])} 없는 함수임")
            return None
        name, line = line[0], line[1:]
        call = self.call[name]
        start = call[0]
        param: List[str] = call[1]
        if "안궁" in "".join(line):
            self.error("안물안궁: 한 줄에서 2번 이상 함수를 호출함")
            return None
        if len(param) != len(line):
            self.error(f"안물안궁: {len(param)}개의 인자가 필요한데, {len(line)}개가 주어짐")
            return None
        new_line: List[int] = []
        for value in line:
            value = self.calc(value)
            if value is None:
                return None
            if type(value) == str:
                value = ord(value)
            new_line.append(value)
        dic: Dict[str, int] = dict(zip(param, new_line))
        self.funcs.append(func(name, start, dic))

    def make_var(self, line: str) -> None:
        line = line.split("~", 1)
        if len(line) == 1:
            name = line[0]
            value = ""
        else:
            name, value = line[0], line[1]
        include = self.check_name(name)
        if include:
            self.error(f"어쩔변수: \"{name}\"{end_letter(name)} 키워드 \"{include}\"{end_letter(include, '을', '를')} 포함함")
            return
        if not name:
            self.error("어쩔변수: 변수 이름이 필요함")
            return
        if name in self.funcs[-1].var or name in self.funcs[-1].var_uni:
            self.error(f"어쩔변수: \"{name}\"{end_letter(name)} 이미 선언됨")
            return
        value = self.calc(value)
        if value is None:
            return
        if type(value) == str:
            value = ord(value)
        self.funcs[-1].var[name] = value

    def assign_var(self, line: str) -> None:
        line = line.split("~", 1)
        if len(line) == 1:
            name = line[0]
            value = ""
        else:
            name, value = line[0], line[1]
        if name not in self.funcs[-1].var:
            self.error(f"어쩔변수: \"{name}\"{end_letter(name)} 선언되지 않음")
            return
        if name in self.funcs[-1].var_uni:
            self.error(f"어쩔변수: \"{name}\"{end_letter(name)} 유니코드 변수임")
            return
        value = self.calc(value)
        if value is None:
            return
        if type(value) == str:
            value = ord(value)
        self.funcs[-1].var[name] = value

    def make_var_uni(self, line: str) -> None:
        line = line.split("~", 1)
        if len(line) == 1:
            name = line[0]
            value = ""
        else:
            name, value = line[0], line[1]
        include = self.check_name(name)
        if include:
            self.error(f"어쩔변수: \"{name}\"{end_letter(name)} 키워드 \"{include}\"{end_letter(include, '을', '를')} 포함함")
            return
        if not name:
            self.error("어쩔변수: 변수 이름이 필요함")
            return
        if name in self.funcs[-1].var or name in self.funcs[-1].var_uni:
            self.error(f"어쩔변수: \"{name}\"{end_letter(name)} 이미 선언됨")
            return
        value = self.calc(value)
        if value is None:
            return
        if type(value) == str:
            value = ord(value)
        self.funcs[-1].var_uni[name] = value

    def assign_var_uni(self, line: str) -> None:
        line = line.split("~", 1)
        if len(line) == 1:
            name = line[0]
            value = ""
        else:
            name, value = line[0], line[1]
        if name not in self.funcs[-1].var_uni:
            self.error(f"어쩔변수: \"{name}\"{end_letter(name)} 선언되지 않음")
            return
        if name in self.funcs[-1].var:
            self.error(f"어쩔변수: \"{name}\"{end_letter(name)} 정수형 변수임")
            return
        value = self.calc(value)
        if value is None:
            return
        if type(value) == str:
            value = ord(value)
        self.funcs[-1].var_uni[name] = value

    def print(self, line: str) -> None:
        value = self.calc(line)
        if value is None:
            return
        print(value, end=("" if type(value) is str and self.file else "\n"))

    def execute_line(self, line: str) -> None:
        if line.strip() == "":
            return
        if not line.startswith(self.keywords):
            self.error("실행놈아: 실행 가능한 구문이 아님")
            return
        for i, keyword in enumerate(self.keywords):
            if line.startswith(keyword):
                line = line[len(keyword):]
                line = "".join(line.split()).strip()
                do = self.exec[i]
                if not do:
                    self.error("실행놈아: 실행 가능한 구문이 아님")
                    return
                do(line)
                return

    def execute(self) -> None:
        print(f"asserlang-python interpreter {self.version}")
        print("by sangchoo1201")
        print(">>> 쿠쿠루삥뽕")
        self.lines.append("쿠쿠루삥뽕")
        self.funcs[-1].cnt = 2
        while True:
            if self.funcs[-1].cnt > len(self.lines):
                line = input("... " if self.writing_func else ">>> ").strip()
                if line == "슉슈슉슉":
                    return
                self.lines.append(line)
            if not self.writing_func or self.lines[self.funcs[-1].cnt-1].startswith("안물"):
                self.execute_line(self.lines[self.funcs[-1].cnt-1])
                if self.stop and len(self.funcs) > 1:
                    self.retn(None)
            self.funcs[-1].cnt += 1
            self.stop = False

    def execute_all(self, lines: List[str], path: Optional[str] = "__default__") -> None:
        if lines[0] not in ("쿠쿠루삥뽕", "ㅋㅋ루삥뽕"):
            print("아무것도 모르죠?: 어쩔랭은 \"쿠쿠루삥뽕\"으로 시작해야 함")
            return
        if lines[-1] != "슉슈슉슉":
            print("아무것도 모르죠?: 어쩔랭은 \"슉슈슉슉\"으로 끝나야 함")
            return
        self.lines = lines[:-1]
        self.funcs[-1].cnt = 2
        self.file = path
        while True:
            if not self.writing_func or self.lines[self.funcs[-1].cnt-1].startswith("안물"):
                self.execute_line(self.lines[self.funcs[-1].cnt-1])
            if self.stop:
                return
            self.funcs[-1].cnt += 1
            if self.funcs[-1].cnt > len(self.lines):
                return

    def execute_file(self, path: str) -> None:
        if not os.path.exists(path):
            print(f"어쩔파일: \"{path}\"에는 파일이 없음")
            return
        with open(path, "r", encoding="utf8") as f:
            lines = f.read().strip().split("\n")
        self.execute_all(lines, path)


if __name__ == "__main__":
    compiler = asserlang()
    arg = sys.argv
    if len(arg) > 1:
        compiler.execute_file(arg[1])
    else:
        compiler.execute()
