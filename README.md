# MasteryPath - Track and Master Your Skills

![MIT License](https://img.shields.io/badge/License-MIT-green.svg)
![Go](https://img.shields.io/badge/Backend-Go_1.21+-00ADD8.svg?logo=go&logoColor=white)
![Expo](https://img.shields.io/badge/Mobile-Expo_50+-000020.svg?logo=expo&logoColor=white)
![MongoDB](https://img.shields.io/badge/Database-MongoDB-47A248.svg?logo=mongodb&logoColor=white)
![React Native](https://img.shields.io/badge/UI-React_Native-61DAFB.svg?logo=react&logoColor=black)

**MasteryPath** is a comprehensive skill tracking application designed to help you monitor, analyze, and master your skills over time. Whether you're learning programming languages, musical instruments, or professional competencies, MasteryPath provides the tools to visualize your progress and identify areas for growth.

## 🚀 The GEM-R Stack

Built with a modern, powerful technology stack:

*   **G**o (Golang): High-performance backend API for fast data processing and business logic.
*   **E**xpo: React Native framework for a seamless, cross-platform mobile experience.
*   **M**ongoDB: Flexible NoSQL database optimized for skill data and progress tracking.
*   **R**eact Native: Component-based UI library for beautiful, native-feeling interfaces.

## ✨ Key Features

*   **📊 Skill Tracking**: Create and monitor multiple skills with detailed progress metrics
*   **📈 Visual Analytics**: Radar charts and growth visualizations to see your development at a glance
*   **🎯 Goal Setting**: Set targets and milestones for each skill you're mastering
*   **📱 Cross-Platform**: Native mobile experience on both iOS and Android
*   **⚡ Real-time Sync**: Instant data synchronization across devices

## 🎯 Project Goals

*   **Empower Learning**: Provide actionable insights into skill development and progress
*   **Data-Driven Growth**: Help users make informed decisions about where to focus their efforts
*   **Intuitive Experience**: Deliver a beautiful, easy-to-use interface for tracking personal development
*   **Scalable Architecture**: Demonstrate modern full-stack development practices with the GEM-R stack

## 🛠️ Installation & Walkthrough

If you'd like to run this project locally to see it in action, follow these steps.

### Prerequisites

*   [Go](https://go.dev/dl/) (v1.21+)
*   [Node.js](https://nodejs.org/) (LTS)
*   [MongoDB](https://www.mongodb.com/try/download/community) (Local or Atlas)

### 1. Environment & Database
Ensure your MongoDB instance is running.
```powershell
# Verify MongoDB service (Windows)
get-service MongoDB
```

### 2. Backend Setup (Go)
The backend service connects to MongoDB and exposes API endpoints.

```bash
cd backend
# 1. Create .env file for configuration
cp .env.example .env
# (Windows Command Prompt: copy .env.example .env)

# 2. Install dependencies
go mod tidy

# 3. Start the server
go run main.go
```

### 3. Frontend Setup (Expo)
The frontend application connects to the Go backend.

```bash
cd frontend
# 1. Configure API URL
cp .env.example .env
# (Windows Command Prompt: copy .env.example .env)

# 2. Install Node dependencies
npm install

# 3. Start the Expo development server
npx expo start
```

**To Run on Your Device:**
*   Download the **Expo Go** app on iOS or Android.
*   Scan the QR code displayed in the terminal.

## 🚀 Going to Production
Ready to deploy? Check out our detailed [Deployment Guide](DEPLOYMENT.md) for instructions on:
*   ☁️ **Database**: Setting up MongoDB Atlas.
*   🐳 **Backend**: Dockerizing the Go API.
*   📱 **Frontend**: Building for App Stores (EAS) and Web (Vercel).

## 🔮 Future Improvements

*   **🔐 User Authentication**: JWT-based authentication for multi-user support
*   **📊 Advanced Analytics**: Trend analysis, skill correlations, and predictive insights
*   **🏆 Achievements System**: Gamification with badges and milestones
*   **👥 Social Features**: Share progress and compare skills with friends
*   **📤 Data Export**: Export your skill data in various formats (CSV, JSON, PDF reports)
*   **🧪 Testing Suite**: Comprehensive unit and integration tests

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

Copyright (c) 2026 Deepak Thomas


