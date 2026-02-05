package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestBuildQueueWorkerCommand_Basic(t *testing.T) {
	plan := ForgeWorkerResourceModel{
		PHPVersion:       types.StringValue("php82"),
		WorkerConnection: types.StringValue("redis"),
		Queue:            types.StringNull(),
		Delay:            types.Int64Value(0),
		Memory:           types.Int64Value(128),
		Sleep:            types.Int64Value(3),
		Timeout:          types.Int64Value(60),
		Tries:            types.Int64Value(0),
		Force:            types.BoolValue(false),
		Daemon:           types.BoolValue(true),
	}

	cmd := buildQueueWorkerCommand(plan)
	expected := "php82 artisan queue:work redis --memory=128 --sleep=3 --timeout=60 --daemon"
	if cmd != expected {
		t.Errorf("expected:\n  %s\ngot:\n  %s", expected, cmd)
	}
}

func TestBuildQueueWorkerCommand_AllOptions(t *testing.T) {
	plan := ForgeWorkerResourceModel{
		PHPVersion:       types.StringValue("php81"),
		WorkerConnection: types.StringValue("sqs"),
		Queue:            types.StringValue("high,default"),
		Delay:            types.Int64Value(5),
		Memory:           types.Int64Value(256),
		Sleep:            types.Int64Value(1),
		Timeout:          types.Int64Value(120),
		Tries:            types.Int64Value(3),
		Force:            types.BoolValue(true),
		Daemon:           types.BoolValue(true),
	}

	cmd := buildQueueWorkerCommand(plan)
	expected := "php81 artisan queue:work sqs --queue=high,default --delay=5 --memory=256 --sleep=1 --timeout=120 --tries=3 --force --daemon"
	if cmd != expected {
		t.Errorf("expected:\n  %s\ngot:\n  %s", expected, cmd)
	}
}

func TestBuildQueueWorkerCommand_DefaultPHP(t *testing.T) {
	plan := ForgeWorkerResourceModel{
		PHPVersion:       types.StringValue(""),
		WorkerConnection: types.StringValue("database"),
		Queue:            types.StringNull(),
		Delay:            types.Int64Value(0),
		Memory:           types.Int64Value(128),
		Sleep:            types.Int64Value(3),
		Timeout:          types.Int64Value(60),
		Tries:            types.Int64Value(0),
		Force:            types.BoolValue(false),
		Daemon:           types.BoolValue(false),
	}

	cmd := buildQueueWorkerCommand(plan)
	expected := "php artisan queue:work database --memory=128 --sleep=3 --timeout=60"
	if cmd != expected {
		t.Errorf("expected:\n  %s\ngot:\n  %s", expected, cmd)
	}
}

func TestParseMemoryFromCommand(t *testing.T) {
	tests := []struct {
		name     string
		command  string
		expected int
		hasError bool
	}{
		{"equals format", "php artisan queue:work --memory=256", 256, false},
		{"space format", "php artisan queue:work --memory 128", 128, false},
		{"with other flags", "php82 artisan queue:work redis --queue=default --memory=512 --sleep=3", 512, false},
		{"no memory flag", "php artisan queue:work redis", 0, true},
		{"empty command", "", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mem, err := parseMemoryFromCommand(tt.command)
			if tt.hasError {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if mem != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, mem)
			}
		})
	}
}

func TestParseCronExpressionIntoModel(t *testing.T) {
	tests := []struct {
		name            string
		cron            string
		expectedMinute  string
		expectedHour    string
		expectedDay     string
		expectedMonth   string
		expectedWeekday string
		hasError        bool
	}{
		{"standard cron", "*/5 * * * *", "*/5", "*", "*", "*", "*", false},
		{"specific time", "30 2 15 6 1", "30", "2", "15", "6", "1", false},
		{"too few parts", "* *", "", "", "", "", "", true},
		{"extra spaces", "  0   12   *   *   0  ", "0", "12", "*", "*", "0", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := &ForgeScheduledJobResourceModel{}
			err := parseCronExpressionIntoModel(tt.cron, model)
			if tt.hasError {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if model.Minute.ValueString() != tt.expectedMinute {
				t.Errorf("minute: expected '%s', got '%s'", tt.expectedMinute, model.Minute.ValueString())
			}
			if model.Hour.ValueString() != tt.expectedHour {
				t.Errorf("hour: expected '%s', got '%s'", tt.expectedHour, model.Hour.ValueString())
			}
			if model.Day.ValueString() != tt.expectedDay {
				t.Errorf("day: expected '%s', got '%s'", tt.expectedDay, model.Day.ValueString())
			}
			if model.Month.ValueString() != tt.expectedMonth {
				t.Errorf("month: expected '%s', got '%s'", tt.expectedMonth, model.Month.ValueString())
			}
			if model.Weekday.ValueString() != tt.expectedWeekday {
				t.Errorf("weekday: expected '%s', got '%s'", tt.expectedWeekday, model.Weekday.ValueString())
			}
		})
	}
}

func TestSplitCompositeID(t *testing.T) {
	tests := []struct {
		name     string
		id       string
		count    int
		expected []string
	}{
		{"two parts", "123:456", 2, []string{"123", "456"}},
		{"three parts", "1:2:3", 3, []string{"1", "2", "3"}},
		{"wrong count", "1:2", 3, nil},
		{"too many parts", "1:2:3", 2, nil},
		{"empty string", "", 2, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := splitCompositeID(tt.id, tt.count)
			if tt.expected == nil {
				if result != nil {
					t.Errorf("expected nil, got %v", result)
				}
				return
			}
			if result == nil {
				t.Fatal("expected non-nil result, got nil")
			}
			for i, v := range tt.expected {
				if result[i] != v {
					t.Errorf("part %d: expected '%s', got '%s'", i, v, result[i])
				}
			}
		})
	}
}
