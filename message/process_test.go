package message

import (
	"sync"
	"testing"
	"time"

	"WaterMark/pkg"
)

func drainInfoChannel() {
	for {
		select {
		case <-Info_Messge_Chan:
		case <-Error_Messge_Chan:
		case <-time.After(10 * time.Millisecond):
			return
		}
	}
}

func resetInfoChannelFlag() {
	msgMtx.Lock()
	info_channel_close_flag = false
	msgMtx.Unlock()
}

func TestSendStartSuccess(t *testing.T) {
	drainInfoChannel()
	resetInfoChannelFlag()
	tests := []struct {
		name string
	}{
		{
			name: "Send start success message",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SendStartSuccess()

			select {
			case msg := <-Info_Messge_Chan:
				if msg != start_success_message {
					t.Errorf("SendStartSuccess() sent %v, want %v", msg, start_success_message)
				}
			case <-time.After(100 * time.Millisecond):
				t.Error("SendStartSuccess() did not send message within timeout")
			}
		})
	}
}

func TestHasSendSuccess(t *testing.T) {
	tests := []struct {
		name     string
		msg      string
		expected bool
	}{
		{
			name:     "Message matches start success",
			msg:      start_success_message,
			expected: true,
		},
		{
			name:     "Message does not match start success",
			msg:      "other message",
			expected: false,
		},
		{
			name:     "Empty message",
			msg:      "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := HasSendSuccess(tt.msg)
			if result != tt.expected {
				t.Errorf("HasSendSuccess() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestSendErrorMsg(t *testing.T) {
	drainInfoChannel()
	tests := []struct {
		name   string
		errStr string
	}{
		{
			name:   "Send error message",
			errStr: "test error",
		},
		{
			name:   "Send empty error message",
			errStr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SendErrorMsg(tt.errStr)

			select {
			case msg := <-Error_Messge_Chan:
				if msg != tt.errStr {
					t.Errorf("SendErrorMsg() sent %v, want %v", msg, tt.errStr)
				}
			case <-time.After(100 * time.Millisecond):
				t.Error("SendErrorMsg() did not send message within timeout")
			}
		})
	}
}

func TestSendErrorOrInfo(t *testing.T) {
	drainInfoChannel()
	resetInfoChannelFlag()
	tests := []struct {
		name     string
		err      pkg.EError
		success  string
		checkErr bool
	}{
		{
			name:     "Send success message",
			err:      pkg.EError{},
			success:  "operation completed",
			checkErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SendErrorOrInfo(tt.err, tt.success)

			if tt.checkErr {
				select {
				case msg := <-Error_Messge_Chan:
					expected := tt.err.String()
					if msg != expected {
						t.Errorf("SendErrorOrInfo() sent error %v, want %v", msg, expected)
					}
				case <-time.After(100 * time.Millisecond):
					t.Error("SendErrorOrInfo() did not send error message within timeout")
				}
			} else {
				select {
				case msg := <-Info_Messge_Chan:
					if msg != tt.success {
						t.Errorf("SendErrorOrInfo() sent info %v, want %v", msg, tt.success)
					}
				case <-time.After(100 * time.Millisecond):
					t.Error("SendErrorOrInfo() did not send info message within timeout")
				}
			}
		})
	}
}

func TestSendInfoMsg(t *testing.T) {
	drainInfoChannel()
	resetInfoChannelFlag()
	tests := []struct {
		name string
		info string
	}{
		{
			name: "Send info message",
			info: "test info",
		},
		{
			name: "Send empty info message",
			info: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SendInfoMsg(tt.info)

			select {
			case msg := <-Info_Messge_Chan:
				if msg != tt.info {
					t.Errorf("SendInfoMsg() sent %v, want %v", msg, tt.info)
				}
			case <-time.After(100 * time.Millisecond):
				t.Error("SendInfoMsg() did not send message within timeout")
			}
		})
	}
}

func TestConcurrentSend(t *testing.T) {
	drainInfoChannel()
	resetInfoChannelFlag()
	tests := []struct {
		name string
	}{
		{
			name: "Concurrent send messages",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			done := make(chan bool)

			for i := 0; i < 10; i++ {
				go func(idx int) {
					SendInfoMsg("info message")
					done <- true
				}(i)
			}

			for i := 0; i < 10; i++ {
				<-done
			}

			for i := 0; i < 10; i++ {
				select {
				case <-Info_Messge_Chan:
				case <-time.After(100 * time.Millisecond):
					t.Error("Concurrent send timeout")
				}
			}
		})
	}
}

func TestChannelCapacity(t *testing.T) {
	drainInfoChannel()
	resetInfoChannelFlag()
	tests := []struct {
		name string
	}{
		{
			name: "Test channel capacity",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for i := 0; i < 100; i++ {
				SendInfoMsg("message")
			}

			count := 0
			for {
				select {
				case <-Info_Messge_Chan:
					count++
				case <-time.After(10 * time.Millisecond):
					if count != 100 {
						t.Errorf("Expected 100 messages, got %d", count)
					}
					return
				}
			}
		})
	}
}

func TestMutexLock(t *testing.T) {
	drainInfoChannel()
	resetInfoChannelFlag()
	tests := []struct {
		name string
	}{
		{
			name: "Test mutex lock",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var wg sync.WaitGroup
			wg.Add(100)

			for i := 0; i < 100; i++ {
				go func() {
					defer wg.Done()
					SendInfoMsg("test message")
				}()
			}

			wg.Wait()

			for i := 0; i < 100; i++ {
				select {
				case <-Info_Messge_Chan:
				case <-time.After(100 * time.Millisecond):
					t.Error("Mutex lock test timeout")
				}
			}
		})
	}
}

func TestSendErrorOrInfoWithError(t *testing.T) {
	drainInfoChannel()
	tests := []struct {
		name    string
		err     pkg.EError
		success string
	}{
		{
			name:    "Send error message",
			err:     pkg.NewErrors(1, "test error"),
			success: "success",
		},
		{
			name:    "Send error with code",
			err:     pkg.NewErrors(100, "file not found"),
			success: "operation completed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SendErrorOrInfo(tt.err, tt.success)

			select {
			case msg := <-Error_Messge_Chan:
				expected := tt.err.String()
				if msg != expected {
					t.Errorf("SendErrorOrInfo() sent error %v, want %v", msg, expected)
				}
			case <-time.After(100 * time.Millisecond):
				t.Error("SendErrorOrInfo() did not send error message within timeout")
			}
		})
	}
}

func TestSendMultipleMessages(t *testing.T) {
	drainInfoChannel()
	resetInfoChannelFlag()
	tests := []struct {
		name string
	}{
		{
			name: "Send multiple messages",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			messages := []string{"msg1", "msg2", "msg3", "msg4", "msg5"}

			for _, msg := range messages {
				SendInfoMsg(msg)
			}

			received := make(map[string]bool)
			for i := 0; i < len(messages); i++ {
				select {
				case msg := <-Info_Messge_Chan:
					received[msg] = true
				case <-time.After(100 * time.Millisecond):
					t.Error("Multiple messages test timeout")
				}
			}

			for _, expectedMsg := range messages {
				if !received[expectedMsg] {
					t.Errorf("Did not receive expected message %v", expectedMsg)
				}
			}
		})
	}
}

func TestSendErrorMultipleMessages(t *testing.T) {
	drainInfoChannel()
	tests := []struct {
		name string
	}{
		{
			name: "Send multiple error messages",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			messages := []string{"error1", "error2", "error3"}

			for _, msg := range messages {
				SendErrorMsg(msg)
			}

			for _, expectedMsg := range messages {
				select {
				case msg := <-Error_Messge_Chan:
					if msg != expectedMsg {
						t.Errorf("Received %v, want %v", msg, expectedMsg)
					}
				case <-time.After(100 * time.Millisecond):
					t.Error("Multiple error messages test timeout")
				}
			}
		})
	}
}
