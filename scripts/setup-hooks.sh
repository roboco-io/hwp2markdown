#!/bin/bash
# Git hooks 설정 스크립트

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
HOOKS_DIR="$PROJECT_ROOT/.githooks"

echo "🔧 Setting up Git hooks..."

# Git hooks 경로 설정
git config core.hooksPath "$HOOKS_DIR"

# 실행 권한 부여
chmod +x "$HOOKS_DIR"/*

echo "✅ Git hooks configured successfully!"
echo ""
echo "Hooks location: $HOOKS_DIR"
echo "Active hooks:"
ls -la "$HOOKS_DIR"
