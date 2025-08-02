## 🛠️ 1foot Resolver Setup Guide

The **Resolver** is written in **Golang** and includes two core services:

- `OrderFulfillmentService`
- `HTLCCreationService`

> ✅ Make sure **Go 1.24+** is installed on your local machine. You can download it from [golang.org/dl](https://golang.org/dl/).

---

### 📦 1. Clone the Repository

```bash
git clone https://github.com/1foot-Labs/1foot-resolver
cd 1foot-resolver
```

### 📐 2. Install Dependencies

```bash
go mod tidy
```

### ▶️ 3. Run the Project

```bash
go run main.go
```

---

> 💡 **Tip:** Use `go run ./cmd/resolver` if your entry point is under a `cmd/` directory.

Let us know if you encounter any issues — happy building! ⚡
