def end(word, yes="은", no="는"):
    if ord("ㄱ") <= ord(word[-1]) <= ord("ㅎ"):
        return yes
    if ord("ㅏ") <= ord(word[-1]) <= ord("ㅣ"):
        return no
    if (ord(word[-1]) - ord("가")) % 28:
        return yes
    else:
        return no

class func:
    def __init__(self, name, start, param):
        self.name = name
        self.start = start
        self.cnt = start
        self.var = param
        self.var_uni = {}

class asserlang:
    def __init__(self):
        self.keywords = ("ㅋ", "ㅎ", "ㅌ",
                         "어쩔", "저쩔", "우짤래미", "저짤래미",
                         "ㅇㅉ", "ㅌㅂ", "안물", "안궁",
                         "무지개반사", "화났쥬?", "킹받쥬?", ";;")
        self.exec = (None, None, None,
                     self.make_var, self.assign_var, self.make_var_uni, self.assign_var_uni,
                     self.print, None, self.make_func, self.call_func,)
                     #self.retn, self.condition, None, self.jump)
        self.call = {}
        self.funcs = [func("__main__", 1, {})]
        self.lines = []
        self.writing_func = False
        self.stop = False
        self.file = ""
        self.return_value = 0

    def error(self, state):
        print("Traceback (most recent call last):")
        for i in self.funcs:
            print("    ", end="")
            if self.file:
                print(f"file \"{self.file}\", ", end="")
            print(f"line {i.cnt}, in {i.name}")
            print(f"      {self.lines[i.cnt-1]}")
        print(state)
        self.stop = True

    def calc(self, value):
        if value:
            return "A"
        else:
            return 0

    def check_name(self, name):
        for i in self.keywords:
            if i in name:
                return i
        return False

    def retn(self, line):
        pass

    def make_func(self, line):
        if line.strip("~") == "" and not self.writing_func:
            if len(self.funcs) == 1:
                self.error("안물: 함수 이름이 필요함")
                return
            self.retn("")
            return
        if line == "":
            self.writing_func = False
            return
        if self.writing_func:
            self.error("안물: 함수 안에서 함수를 작성함")
            del self.call[self.writing_func]
            self.writing_func = False
            return
        line = line.split("~")
        include = self.check_name(line[0])
        if include:
            self.error(f"안물: 함수 \"{line[0]}\"{end(line[0])} 키워드 \"{include}\"{end(include, '을', '를')} 포함함")
            return
        names = []
        for name in line[1:]:
            include = self.check_name(name)
            if include:
                self.error(f"안물: 매개변수 \"{name}\"{end(name)} 키워드 \"{include}\"{end(include, '을', '를')} 포함함")
                return
            if name in names:
                self.error(f"안물: 매개변수 \"{name}\"{end(name)} 다른 매개변수와 겹침")
                return
            names.append(name)
        self.call[line[0]] = (self.funcs[-1].cnt, tuple(names))
        self.writing_func = line[0]

    def call_func(self, line):
        if line.strip("~") == "":
            self.error("안물: 함수 이름이 필요함")
            return None
        line = line.split("~")
        if line[0] not in self.call:
            self.error(f"안물: \"{line[0]}\"{end(line[0])} 없는 함수임")
            return None


    def make_var(self, line):
        line = line.split("~", 1)
        if len(line) == 1:
            name = line[0]
            value = 0
        else:
            name, value = line
        include = self.check_name(name)
        if include:
            self.error(f"어쩔변수: \"{name}\"{end(name)} 키워드 \"{include}\"{end(include, '을', '를')} 포함함")
            return
        if not name:
            self.error("어쩔변수: 변수 이름이 필요함")
            return
        if name in self.funcs[-1].var or name in self.funcs[-1].var_uni:
            self.error(f"어쩔변수: \"{name}\"{end(name)} 이미 선언됨")
            return
        value = self.calc(value)
        if value is None:
            return
        self.funcs[-1].var[name] = ord(value) if type(value) == str else value

    def assign_var(self, line):
        line = line.split("~", 1)
        if len(line) == 1:
            name = line[0]
            value = 0
        else:
            name, value = line
        if name not in self.funcs[-1].var:
            self.error(f"어쩔변수: \"{name}\"{end(name)} 선언되지 않음")
            return
        if name in self.funcs[-1].var_uni:
            self.error(f"어쩔변수: \"{name}\"{end(name)} 유니코드 변수임")
            return
        value = self.calc(value)
        if value is None:
            return
        if type(value) == str:
            value = ord(value)
        self.funcs[-1].var[name] = ord(value) if type(value) == str else value

    def make_var_uni(self, line):
        line = line.split("~", 1)
        if len(line) == 1:
            name = line[0]
            value = 0
        else:
            name, value = line
        include = self.check_name(name)
        if include:
            self.error(f"어쩔변수: \"{name}\"{end(name)} 키워드 \"{include}\"{end(include, '을', '를')} 포함함")
            return
        if not name:
            self.error("어쩔변수: 변수 이름이 필요함")
            return
        if name in self.funcs[-1].var or name in self.funcs[-1].var_uni:
            self.error(f"어쩔변수: \"{name}\"{end(name)} 이미 선언됨")
            return
        value = self.calc(value)
        if value is None:
            return
        if type(value) == str:
            value = ord(value)
        self.funcs[-1].var_uni[name] = ord(value) if type(value) == str else value

    def assign_var_uni(self, line):
        line = line.split("~", 1)
        if len(line) == 1:
            name = line[0]
            value = 0
        else:
            name, value = line
        if name not in self.funcs[-1].var_uni:
            self.error(f"어쩔변수: \"{name}\"{end(name)} 선언되지 않음")
            return
        if name in self.funcs[-1].var:
            self.error(f"어쩔변수: \"{name}\"{end(name)} 정수형 변수임")
            return
        value = self.calc(value)
        if value is None:
            return
        if type(value) == str:
            value = ord(value)
        self.funcs[-1].var_uni[name] = ord(value) if type(value) == str else value

    def print(self, line):
        value = self.calc(line)
        if value is None:
            return
        print(value)

    def execute_line(self, line: str):
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

    def execute(self):
        print("asserlang-python interpreter v1.0")
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
            self.funcs[-1].cnt += 1


if __name__ == "__main__":
    compiler = asserlang()
    compiler.execute()
