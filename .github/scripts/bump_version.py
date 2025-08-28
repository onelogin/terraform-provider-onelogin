import re
import sys

def bump_version(version, bump_type):
    major, minor, patch = map(int, version.split("."))
    if bump_type == "major":
        major += 1
        minor = 0
        patch = 0
    elif bump_type == "minor":
        minor += 1
        patch = 0
    elif bump_type == "patch":
        patch += 1
    else:
        raise ValueError("Invalid bump type")
    return f"{major}.{minor}.{patch}"

def update_gnumakefile(path, new_version):
    with open(path, "r") as f:
        content = f.read()
    content = re.sub(r'VERSION=[0-9]+\.[0-9]+\.[0-9]+', f'VERSION={new_version}', content)
    with open(path, "w") as f:
        f.write(content)

def update_readme(path, new_version):
    with open(path, "r") as f:
        content = f.read()
    content = re.sub(
        r'version = "[0-9]+\.[0-9]+\.[0-9]+"',
        f'version = "{new_version}"',
        content
    )
    with open(path, "w") as f:
        f.write(content)

def extract_version_gnumakefile(path):
    with open(path, "r") as f:
        for line in f:
            m = re.match(r'VERSION=([0-9]+\.[0-9]+\.[0-9]+)', line)
            if m:
                return m.group(1)
    raise ValueError("VERSION not found in GNUmakefile")

if __name__ == "__main__":
    bump_type = sys.argv[1]
    gnumakefile = "GNUmakefile"
    readme = "README.md"

    # Extract current version
    current_version = extract_version_gnumakefile(gnumakefile)
    new_version = bump_version(current_version, bump_type)

    update_gnumakefile(gnumakefile, new_version)
    update_readme(readme, new_version)

    # Output for GitHub Actions
    print(f"::set-output name=new_version::{new_version}")
