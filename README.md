# Genesis Enuma Elish - Backend API

Genesis Enuma Elish is a backend API for a Learning Management System (LMS) built using Go and the Gin framework. This system supports school management, students, teachers, classes, subjects, exams, etc.

## 🏗️ Architecture

### Tech Stack
- **Language**: Go 1.21+
- **Framework**: Gin Web Framework
- **Database**: PostgreSQL
- **Cache**: Redis
- **File Storage**: Cloudinary
- **Validation**: go-playground/validator
- **Migration**: Custom migration system
- **Telemetry**: OpenTelemetry

### Project Structure
```
enki/
├── api/                    # API initialization and routing
├── cmd/
│   ├── api/               # Main application entry point
│   └── migration/         # Database migration tool
├── config/                # Configuration management
├── infra/                 # Infrastructure setup (DB, Redis, etc.)
├── internal/              # Internal modules
│   ├── auth/             # Authentication & authorization
│   ├── class/            # Class management
│   ├── exam/             # Exam system
│   ├── ppdb/             # Student admission system
│   ├── question/         # Question bank
│   ├── school/           # School management
│   ├── storage/          # File storage
│   ├── student/          # Student management
│   ├── subject/          # Subject management
│   └── teacher/          # Teacher management
└── pkg/                   # Shared packages
    ├── error/            # Error handling
    ├── http/             # HTTP utilities
    ├── middleware/       # Custom middleware
    └── migration/        # Migration utilities
```

### Module
Each module follows a pattern:
- **Handler**: HTTP request handling and response
- **Service**: Business logic and validation
- **Repository**: Data access layer

## 🚀 Setup & Installation

### Prerequisites
- Go 1.21 or newer
- PostgreSQL 13+
- Redis 6+
- Git

### Environment Setup

1. **Clone repository**
```bash
git clone <repository-url>
cd enki
```

2. **Install dependencies**
```bash
go mod download
```

3. **Setup configuration**
```bash
cp config.example.json config.json
```

4. **Edit configuration** (`config.json`)
```json
{
  "http": {
    "host": "0.0.0.0",
    "port": "8000",
    "read_timeout": 30,
    "write_timeout": 30,
    "frontend_host": "http://localhost:5173"
  },
  "postgres": {
    "host": "localhost",
    "port": "5432",
    "username": "postgres",
    "password": "password",
    "database": "genesis",
    "ssl_mode": "disable"
  },
  "redis": {
    "host": "localhost",
    "port": "6379"
  },
  "jwt": {
    "secret": "your-jwt-secret",
    "duration": 60
  },
  "cloudinary": {
    "cloud_name": "your_cloud_name",
    "api_key": "your_api_key",
    "api_secret": "your_api_secret",
    "folder": "genesis"
  }
}
```

### Database Setup

1. **Run migrations**
```bash
go run cmd/migration/main.go migrate up
```

2. **Run seeder (optional)**
```bash
go run cmd/migration/main.go seed up
```

### Running Application

1. **Development mode**
```bash
go run cmd/api/main.go api
```

2. **Build production**
```bash
go build -o genesis cmd/api/main.go
./genesis api
```

API will run on `http://localhost:8000`

## 📚 API Documentation

### Base URL
```
http://localhost:8000/api/v1
```

### Authentication
All endpoints (except auth) require JWT token in header:
```
Authorization: Bearer <jwt_token>
```

### Core Modules

#### 🔐 Authentication (`/auth`)
- `POST /auth/register` - User registration
- `POST /auth/login` - User login
- `POST /auth/register/verify-email` - Email verification
- `POST /auth/forgot-password` - Forgot password
- `POST /auth/refresh-token` - Refresh JWT token
- `GET /auth/me` - Get current user info
- `PUT /auth/me` - Update user profile

#### 🏫 School Management (`/school`)
- `POST /school` - Create school
- `GET /school` - List schools
- `GET /school/:school_id` - Get school details
- `PUT /school/:school_id` - Update school
- `DELETE /school/:school_id` - Delete school
- `GET /school/statistic` - School statistics
- `GET /school/:school_id/switch` - Switch active school

#### 👨‍🎓 Student Management (`/student`)
- `GET /student` - List students
- `GET /student/:student_id` - Get student details
- `POST /student/invite` - Invite student
- `POST /student/invite/verify` - Verify student email
- `POST /student/invite/complete` - Complete student registration
- `PUT /student/class` - Update student class
- `DELETE /student/:student_id` - Delete student

#### 👨‍🏫 Teacher Management (`/teacher`)
- `GET /teacher` - List teachers
- `GET /teacher/:teacher_id` - Get teacher details
- `POST /teacher/invite` - Invite teacher
- `POST /teacher/invite/verify` - Verify teacher email
- `POST /teacher/invite/complete` - Complete teacher registration
- `GET /teacher/statistic` - Teacher statistics
- `GET /teacher/:teacher_id/subjects` - Get teacher subjects
- `GET /teacher/:teacher_id/classes` - Get teacher classes
- `PUT /teacher/class` - Update teacher class
- `DELETE /teacher/:teacher_id` - Delete teacher

#### 🏛️ Class Management (`/class`)
- `POST /class` - Create class
- `GET /class` - List classes
- `GET /class/:class_id` - Get class details
- `PUT /class/:class_id` - Update class
- `DELETE /class/:class_id` - Delete class
- `POST /class/add-students` - Add students to class
- `POST /class/assign-teachers` - Assign teachers to class
- `POST /class/add-subjects` - Add subjects to class
- `GET /class/:class_id/students` - Get class students
- `GET /class/:class_id/teachers` - Get class teachers
- `GET /class/:class_id/subjects` - Get class subjects

#### 📚 Subject Management (`/subject`)
- `POST /subject` - Create subject
- `GET /subject` - List subjects
- `GET /subject/:subject_id` - Get subject details
- `PUT /subject/:subject_id` - Update subject
- `DELETE /subject/:subject_id` - Delete subject
- `POST /subject/assign-teachers` - Assign teachers to subject
- `GET /subject/:subject_id/teachers` - Get subject teachers
- `PUT /subject/class` - Update subject class

#### 📝 Exam Management (`/exam`)
- `POST /exam` - Create exam
- `GET /exam` - List exams
- `GET /exam/:exam_id` - Get exam details
- `PUT /exam/:exam_id` - Update exam
- `DELETE /exam/:exam_id` - Delete exam
- `POST /exam/assign` - Assign exam to class
- `POST /exam/grade` - Grade exam
- `GET /exam/:exam_id/students` - Get exam students

#### 📝 Student Exam (`/student/exam`)
- `GET /student/exam` - Get student exams
- `GET /student/exam/:exam_id` - Get student exam details
- `POST /student/exam/submit` - Submit exam answers

#### ❓ Question Management (`/question`)
- `POST /question` - Create question
- `GET /question` - List questions
- `GET /question/:question_id` - Get question details
- `PUT /question/:question_id` - Update question
- `DELETE /question/:question_id` - Delete question

#### 🎓 PPDB Management (`/ppdb`)
- `POST /ppdb` - Create PPDB program
- `GET /ppdb` - List PPDB programs
- `GET /ppdb/:ppdb_id` - Get PPDB details
- `PUT /ppdb/:ppdb_id` - Update PPDB
- `DELETE /ppdb/:ppdb_id` - Delete PPDB
- `POST /ppdb/register` - Register for PPDB
- `GET /ppdb/registrants` - Get PPDB registrants
- `POST /ppdb/select` - Select PPDB students

#### 💾 Storage Management (`/storage`)
- `POST /storage/image` - Upload image
- `POST /storage/video` - Upload video
- `POST /storage/document` - Upload document
- `DELETE /storage/file` - Delete file
- `GET /storage/file/:publicId` - Get file info
- `GET /storage/serve/:publicId` - Serve file
- `GET /storage/history` - Get storage history

### Response Format
All responses use standard format:

**Success Response:**
```json
{
  "code": 200,
  "message": "Success message",
  "data": {...},
  "meta": {
    "page": 1,
    "limit": 10,
    "total": 100
  }
}
```

**Error Response:**
```json
{
  "code": 400,
  "message": "Error message",
  "errors": {
    "field_name": ["Validation error message"]
  }
}
```

## 🗃️ Database Schema

### Core Tables
- `users` - User accounts (students, teachers, admins)
- `schools` - School information
- `classes` - Class/grade information
- `subjects` - Subject/course information
- `exams` - Exam information
- `questions` - Question bank
- `ppdb` - Student admission programs
- `ppdb_student` - Student admissions

### Relationship Tables
- `user_school_role` - User roles in schools
- `class_students` - Student-class relationships
- `class_teachers` - Teacher-class relationships
- `subject_teachers` - Teacher-subject relationships
- `exam_classes` - Exam-class assignments
- `student_exam_answers` - Student exam submissions

## 🔧 Development

### Code Standards
- Use `gofmt` for formatting
- Follow Go naming conventions
- Write unit tests for business logic
- Use meaningful commit messages

### Testing
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific test
go test ./internal/auth/test/
```

### Migration
```bash
# Create new migration
go run cmd/migration/main.go create <migration_name>

# Run migrations
go run cmd/migration/main.go migrate up

# Rollback migration
go run cmd/migration/main.go migrate down

# Check migration status
go run cmd/migration/main.go migrate status
```

### Adding New Module

1. **Create module structure:**
```
internal/newmodule/
├── handler/
├── service/
├── repository/
└── newmodule.go
```

2. **Implement interfaces:**
- Repository interface with implementation
- Service interface with business logic
- Handler with HTTP endpoints

3. **Register in `api/api.go`:**
```go
newmodule.New(api.config, api.infra, api.Engine, validate).Init()
```

## 🚨 Error Handling

### Custom Errors
```go
// pkg/error/error.go
var (
    ErrNotFound = New(404, "Resource not found")
    ErrUnauthorized = New(401, "Unauthorized")
    ErrBadRequest = New(400, "Bad request")
)
```

### Validation Errors
Using `go-playground/validator` with custom readable error messages.

## 📊 Monitoring & Telemetry

### OpenTelemetry
- Distributed tracing
- Metrics collection
- Performance monitoring

### Logging
Using `zerolog` for structured logging:
```go
log.Info().Str("user_id", userID).Msg("User logged in")
log.Error().Err(err).Msg("Database error")
```

## 🚀 Deployment

### Docker
```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o genesis cmd/api/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/genesis .
COPY --from=builder /app/config.json .
CMD ["./genesis", "api"]
```

### Environment Variables
```bash
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=password
export DB_NAME=genesis
export REDIS_HOST=localhost
export REDIS_PORT=6379
export JWT_SECRET=your-secret
```

## 🤝 Contributing

1. Fork repository
2. Create feature branch (`git checkout -b feature/new-feature`)
3. Commit changes (`git commit -m 'Add new feature'`)
4. Push to branch (`git push origin feature/new-feature`)
5. Create Pull Request

### Code Review Checklist
- [ ] Tests passing
- [ ] Code formatted with `gofmt`
- [ ] Documentation updated
- [ ] Migration files included
- [ ] Security considerations addressed

## 📄 License
This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 Troubleshooting

### Common Issues

**Database Connection Error:**
```bash
# Check PostgreSQL service
sudo systemctl status postgresql

# Check connection
psql -h localhost -U postgres -d genesis
```

**Redis Connection Error:**
```bash
# Check Redis service
sudo systemctl status redis

# Test connection
redis-cli ping
```

**Migration Errors:**
```bash
# Reset migrations (CAUTION: Data loss)
go run cmd/migration/main.go migrate reset

# Check migration status
go run cmd/migration/main.go migrate status
```

### Debug Mode
```bash
# Run with debug logging
GIN_MODE=debug go run cmd/api/main.go api
```

## 📞 Support

For questions and support:
- Create issue in repository
- Email: ne.nekonako@gmail.com
- Documentation: [Link to docs]
