import tkinter as tk
import subprocess

class GUI:
    def __init__(self, master):
        self.master = master
        master.title("Python GUI with Golang logic")

        self.label = tk.Label(master, text="Enter a number:")
        self.label.pack()

        self.entry = tk.Entry(master)
        self.entry.pack()

        self.button = tk.Button(master, text="Calculate", command=self.calculate)
        self.button.pack()

        self.result_label = tk.Label(master, text="")
        self.result_label.pack()

    def calculate(self):
        number = self.entry.get()
        result = subprocess.check_output(["./calc.exe", number]).decode("utf-8")
        self.result_label.config(text=f"The result is: {result}")

if __name__ == "__main__":
    root = tk.Tk()
    gui = GUI(root)
    root.mainloop()
