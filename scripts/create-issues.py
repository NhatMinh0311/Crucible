#!/usr/bin/env python3
"""Tạo nhãn, milestone, issue và Project board trên GitHub từ .github/backlog.yaml.

Cách dùng:
    GITHUB_TOKEN=$(gh auth token) python3 scripts/create-issues.py <chủ sở hữu>/<kho mã> \
        --m1 <tài khoản> --m2 <tài khoản> --m3 <tài khoản> [--project <tên board>]

Token cần quyền repo; với --project cần thêm quyền project (gh auth login -s project).
Chạy lại an toàn: nhãn, milestone, issue đã tồn tại (theo tiêu đề) được bỏ qua; issue có done: true được đóng;
board trùng tên đã liên kết với kho mã được dùng lại.
"""
import argparse
import json
import os
import sys
import urllib.error
import urllib.request
from pathlib import Path

import yaml

API = "https://api.github.com"
LABEL_COLORS = {"m1": "1d76db", "m2": "0e8a16", "m3": "fbca04", "platform": "5319e7"}


def call(method, path, token, body=None):
    req = urllib.request.Request(API + path, method=method,
                                 data=json.dumps(body).encode() if body is not None else None,
                                 headers={"Authorization": f"Bearer {token}",
                                          "Accept": "application/vnd.github+json",
                                          "X-GitHub-Api-Version": "2022-11-28"})
    try:
        with urllib.request.urlopen(req) as resp:
            return json.loads(resp.read() or "null")
    except urllib.error.HTTPError as err:
        sys.exit(f"{method} {path}: {err.code} {err.read().decode()}")


def graphql(token, query, **variables):
    data = call("POST", "/graphql", token, {"query": query, "variables": variables})
    if data.get("errors"):
        sys.exit(f"GraphQL: {data['errors']}")
    return data["data"]


def project_board(token, repo, title):
    """Trả về id của Project (v2) tên title liên kết với kho mã; tạo mới nếu chưa có."""
    owner, name = repo.split("/")
    r = graphql(token, """query($owner: String!, $name: String!) {
        repository(owner: $owner, name: $name) { id owner { id } projectsV2(first: 50) { nodes { id title url } } } }""",
                owner=owner, name=name)["repository"]
    for p in r["projectsV2"]["nodes"]:
        if p["title"] == title:
            return p["id"]
    p = graphql(token, """mutation($owner: ID!, $repo: ID!, $title: String!) {
        createProjectV2(input: {ownerId: $owner, repositoryId: $repo, title: $title}) { projectV2 { id url } } }""",
                owner=r["owner"]["id"], repo=r["id"], title=title)["createProjectV2"]["projectV2"]
    print("board", p["url"])
    return p["id"]


def paged(path, token):
    page, items = 1, []
    while True:
        batch = call("GET", f"{path}{'&' if '?' in path else '?'}per_page=100&page={page}", token)
        items += batch
        if len(batch) < 100:
            return items
        page += 1


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("repo")
    parser.add_argument("--project", help="tên Project board; bỏ qua thì không tạo board")
    for m in ("m1", "m2", "m3"):
        parser.add_argument(f"--{m}", help=f"tài khoản GitHub phụ trách {m.upper()}")
    args = parser.parse_args()
    token = os.environ.get("GITHUB_TOKEN") or sys.exit("Cần biến môi trường GITHUB_TOKEN")
    backlog = yaml.safe_load(Path(__file__).resolve().parent.parent.joinpath(".github/backlog.yaml").read_text())
    repo = f"/repos/{args.repo}"

    existing_labels = {l["name"] for l in paged(f"{repo}/labels", token)}
    for name in sorted({l for i in backlog["issues"] for l in i["labels"]}):
        if name not in existing_labels:
            call("POST", f"{repo}/labels", token, {"name": name, "color": LABEL_COLORS.get(name, "ededed")})
            print("nhãn", name)

    milestones = {m["title"]: m["number"] for m in paged(f"{repo}/milestones?state=all", token)}
    for key, desc in backlog["milestones"].items():
        if key not in milestones:
            milestones[key] = call("POST", f"{repo}/milestones", token, {"title": key, "description": desc})["number"]
            print("milestone", key)

    issues = {i["title"]: i for i in paged(f"{repo}/issues?state=all", token) if "pull_request" not in i}
    owners = {"m1": args.m1, "m2": args.m2, "m3": args.m3}
    for item in backlog["issues"]:
        if item["title"] in issues:
            continue
        body = {"title": item["title"], "labels": item["labels"], "milestone": milestones[item["milestone"]],
                "body": f"Tham chiếu: {item['refs']}\n\nTiêu chí hoàn thành: CR-DOC-06 mục 6."}
        if owners.get(item["owner"]):
            body["assignees"] = [owners[item["owner"]]]
        issues[item["title"]] = created = call("POST", f"{repo}/issues", token, body)
        print(f"#{created['number']} {item['title']}")

    for item in backlog["issues"]:
        issue = issues[item["title"]]
        if item.get("done") and issue["state"] == "open":
            issues[item["title"]] = call("PATCH", f"{repo}/issues/{issue['number']}", token,
                                         {"state": "closed", "state_reason": "completed"})
            print(f"đóng #{issue['number']}")

    if args.project:
        board = project_board(token, args.repo, args.project)
        for item in backlog["issues"]:  # thêm item đã có trên board không tạo bản trùng
            graphql(token, """mutation($board: ID!, $issue: ID!) {
                addProjectV2ItemById(input: {projectId: $board, contentId: $issue}) { item { id } } }""",
                    board=board, issue=issues[item["title"]]["node_id"])
        print("đã thêm", len(backlog["issues"]), "issue vào board")


if __name__ == "__main__":
    main()
