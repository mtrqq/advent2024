import argparse


def cmdline_arguments():
    parser = argparse.ArgumentParser("advent2024_day1")
    parser.add_argument("--input", default="input.txt")
    return parser.parse_args()


def calculateDistance(path: str):
    with open(path) as fd:
        lines = fd.readlines()

    left, right = [], []
    for index, line in enumerate(lines):
        tokens = line.split(" ", maxsplit=1)

        if len(tokens) != 2:
            raise RuntimeError(f"Invalid number of numbers on the line#{index}")

        left.append(int(tokens[0].strip()))
        right.append(int(tokens[1].strip()))

    left.sort()
    right.sort()

    distance_total = 0
    for left_num, right_num in zip(left, right):
        distance_total += abs(left_num - right_num)

    return distance_total


def main():
    args = cmdline_arguments()
    print("Lists distance:", calculateDistance(args.input))


if __name__ == "__main__":
    main()
