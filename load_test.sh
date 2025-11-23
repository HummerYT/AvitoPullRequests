#!/bin/bash

echo "Starting load test"

echo "1. Health endpoint (30 seconds, 10 threads, 10 connections):"
wrk -t10 -c10 -d30s http://localhost:8080/health

echo ""
echo "2. Stats endpoint (20 seconds, 5 threads, 5 connections):"
wrk -t5 -c5 -d20s http://localhost:8080/stats

echo ""
echo "3. Team get endpoint (15 seconds, 3 threads, 3 connections):"
wrk -t3 -c3 -d15s "http://localhost:8080/team/get?team_name=backend"

echo ""
echo "Load test completed!"