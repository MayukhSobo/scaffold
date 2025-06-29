# 🚀 GitHub Actions Workflows

This directory contains comprehensive CI/CD workflows for the Scaffold project, designed to work seamlessly with the existing Task-based build system.

## 📋 Workflow Overview

### 🔄 CI Pipeline (`ci.yml`)
**Triggers:** Push to `main`/`develop`, Pull Requests to `main`

**Jobs:**
- **Test**: Runs comprehensive test suite with race detection and coverage
- **Lint**: Code quality checks using golangci-lint and vulnerability scanning
- **Build**: Cross-platform builds for Linux, macOS, and Windows
- **Security**: Security scanning with Gosec

**Features:**
- ✅ Leverages existing Task commands
- ✅ Go module caching for faster builds
- ✅ Codecov integration for coverage reporting
- ✅ UPX binary compression
- ✅ Build artifact uploads

### 🏷️ Release Pipeline (`release.yml`)
**Triggers:** Git tags matching `v*`

**Jobs:**
- **Build**: Multi-platform release binaries
- **Docker**: Multi-arch container images pushed to GHCR
- **Release**: GitHub Release creation with changelog
- **Notify**: Success/failure notifications

**Features:**
- ✅ Automated changelog generation
- ✅ SHA256 checksums for binaries
- ✅ Docker images with semantic versioning
- ✅ GitHub Container Registry integration
- ✅ Support for pre-releases

### 🐳 Docker Pipeline (`docker.yml`)
**Triggers:** Push to `main`, Pull Requests

**Jobs:**
- **Build**: Multi-arch Docker image builds
- **Test**: Container functionality testing
- **Security**: Trivy vulnerability scanning

**Features:**
- ✅ Multi-platform builds (amd64, arm64)
- ✅ GitHub Container Registry
- ✅ Docker layer caching
- ✅ Security scanning with SARIF reports
- ✅ Container testing

### 📦 Dependencies (`dependencies.yml`)
**Triggers:** Weekly schedule (Mondays), Manual dispatch

**Jobs:**
- **Security Audit**: Vulnerability and outdated dependency checks
- **Update Dependencies**: Automated patch-level updates via PR
- **Go Version Check**: Monitors for new Go releases
- **License Check**: Dependency license compliance

**Features:**
- ✅ Automated security audits
- ✅ Safe patch-level updates
- ✅ Go version monitoring
- ✅ License compliance tracking
- ✅ Automated PR creation

### 🔍 CodeQL Analysis (`codeql.yml`)
**Triggers:** Push to `main`/`develop`, PRs, Weekly schedule

**Jobs:**
- **Analyze**: Advanced security analysis using GitHub CodeQL

**Features:**
- ✅ Security-extended queries
- ✅ Quality analysis
- ✅ SARIF report integration
- ✅ GitHub Security tab integration

## 🛠️ Workflow Integration

### Task Command Usage
All workflows leverage your existing Task-based build system:

```yaml
# Testing
- run: task test:test
- run: task test:test:race
- run: task test:test:coverage

# Quality
- run: task quality:check
- run: task quality:lint

# Building
- run: task build:build:release:linux
- run: task build:build:release:darwin
- run: task build:build:release:windows

# Dependencies
- run: task deps:install
- run: task deps:vulnerabilities
- run: task deps:outdated
- run: task deps:patch
```

### Caching Strategy
- **Go Modules**: `~/.cache/go-build` and `~/go/pkg/mod`
- **Docker Layers**: GitHub Actions cache
- **Artifacts**: 30-day retention for builds, 1-day for releases

## 🔧 Configuration

### Environment Variables
```yaml
GO_VERSION: '1.22'          # Go version across all workflows
REGISTRY: ghcr.io           # Container registry
```

### Required Secrets
- `GITHUB_TOKEN`: Automatically provided by GitHub
- No additional secrets required for basic functionality

### Optional Enhancements
- **Codecov**: Add `CODECOV_TOKEN` for private repos
- **Slack/Teams**: Add webhook URLs for notifications
- **Custom Registry**: Update `REGISTRY` environment variable

## 📊 Monitoring & Reporting

### Build Artifacts
- **Binaries**: Cross-platform release builds
- **Coverage**: HTML and text coverage reports
- **Security**: SARIF reports for vulnerability scanning
- **Licenses**: CSV reports for dependency licenses

### GitHub Integrations
- **Security Tab**: CodeQL and Trivy findings
- **Actions Tab**: Workflow run history and logs
- **Releases**: Automated releases with binaries
- **Packages**: Container images in GHCR

## 🚀 Getting Started

1. **Push to main**: Triggers CI and Docker workflows
2. **Create PR**: Triggers CI workflow for validation
3. **Create tag**: `git tag v1.0.0 && git push --tags` triggers release
4. **Manual run**: Use "Run workflow" button for dependencies check

### Example Release Process
```bash
# Create and push a release tag
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0

# This triggers:
# ✅ Cross-platform binary builds
# ✅ Docker image creation and push
# ✅ GitHub Release with changelog
# ✅ Artifact uploads
```

## 🔄 Maintenance

### Weekly Automated Tasks
- **Monday 9 AM UTC**: Dependency security audit
- **Monday 2:30 AM UTC**: CodeQL security analysis

### Manual Maintenance
- Update `GO_VERSION` when new Go releases are available
- Review and merge automated dependency PRs
- Monitor security findings in GitHub Security tab

## 📈 Benefits

- **🔄 Fully Automated**: From code push to release
- **🛡️ Security First**: Multiple security scanning layers
- **📦 Multi-Platform**: Builds for all major platforms
- **🐳 Container Ready**: Docker images with multi-arch support
- **📊 Observable**: Comprehensive reporting and monitoring
- **🚀 Fast**: Efficient caching and parallel execution
- **🔧 Maintainable**: Leverages existing Task system 