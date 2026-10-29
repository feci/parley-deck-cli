# Codex native source for round 2

This source material arrived after codex-1 filed round 1. Peers must assess it in cross-review.
The installed CLI reports version 0.165.0. These files use the upstream `rust-v0.165.0` tag.
Tag object: `e374836b2728d9b7aa981c2902649e640a2edaa8`.
No native exhaustion was induced. These excerpts do not establish a captured failure or binary equivalence.
The current adapter launch does not pass `--json`. These source paths describe that explicit output mode.

## codex-rs/exec/src/exec_events.rs
Source: https://github.com/openai/codex/blob/rust-v0.165.0/codex-rs/exec/src/exec_events.rs
Full fetched file SHA256: `dafa872d7e86a099e56e28a329dcb9c03db90ed768c3b88cca8c91d46dc1d0e5`.

Lines 19–24:
```rust
    /// Emitted when a turn is completed. Typically right after the assistant's response.
    #[serde(rename = "turn.completed")]
    TurnCompleted(TurnCompletedEvent),
    /// Indicates that a turn failed with an error.
    #[serde(rename = "turn.failed")]
    TurnFailed(TurnFailedEvent),
```

Lines 54–57:
```rust
#[derive(Debug, Clone, Serialize, Deserialize, PartialEq, TS)]
pub struct TurnFailedEvent {
    pub error: ThreadErrorEvent,
}
```

Lines 90–94:
```rust
/// Fatal error emitted by the stream.
#[derive(Debug, Clone, Serialize, Deserialize, PartialEq, TS)]
pub struct ThreadErrorEvent {
    pub message: String,
}
```

## codex-rs/exec/src/event_processor_with_jsonl_output.rs
Source: https://github.com/openai/codex/blob/rust-v0.165.0/codex-rs/exec/src/event_processor_with_jsonl_output.rs
Full fetched file SHA256: `bc296541b385f9b4ca2c28cc638e25a8d6b1332b2cd67f5a03eaa97a3d612b19`.

Lines 538–557:
```rust
                    TurnStatus::Failed => {
                        self.final_message = None;
                        self.emit_final_message_on_shutdown = false;
                        let error = notification
                            .turn
                            .error
                            .map(|error| ThreadErrorEvent {
                                message: match error.additional_details {
                                    Some(details) if !details.is_empty() => {
                                        format!("{} ({details})", error.message)
                                    }
                                    _ => error.message,
                                },
                            })
                            .or_else(|| self.last_critical_error.clone())
                            .unwrap_or_else(|| ThreadErrorEvent {
                                message: "turn failed".to_string(),
                            });
                        events.push(ThreadEvent::TurnFailed(TurnFailedEvent { error }));
                        CodexStatus::InitiateShutdown
```

## codex-rs/exec/src/lib.rs
Source: https://github.com/openai/codex/blob/rust-v0.165.0/codex-rs/exec/src/lib.rs
Full fetched file SHA256: `4229c67ec7a4d9aa21f26c500f1206660f7ccb3a6c047843a3fa7295938a0848`.

Lines 1673–1691:
```rust
fn should_process_notification(
    notification: &ServerNotification,
    thread_id: &str,
    turn_id: &str,
) -> bool {
    match notification {
        ServerNotification::ConfigWarning(_) | ServerNotification::DeprecationNotice(_) => true,
        // TODO(anp) resolve duplicate startup warnings
        ServerNotification::Warning(notification) => notification
            .thread_id
            .as_deref()
            .is_none_or(|candidate| candidate == thread_id),
        ServerNotification::AuthRecoveryStarted(notification)
        | ServerNotification::AuthRecoveryCompleted(notification) => {
            notification.thread_id == thread_id && notification.turn_id == turn_id
        }
        ServerNotification::Error(notification) => {
            notification.thread_id == thread_id && notification.turn_id == turn_id
        }
```

Lines 1719–1724:
```rust
            notification.thread_id == thread_id && notification.turn_id == turn_id
        }
        ServerNotification::TurnCompleted(notification) => {
            notification.thread_id == thread_id && notification.turn.id == turn_id
        }
        ServerNotification::TurnDiffUpdated(notification) => {
```

## codex-rs/protocol/src/error.rs
Source: https://github.com/openai/codex/blob/rust-v0.165.0/codex-rs/protocol/src/error.rs
Full fetched file SHA256: `73252b4f22850de45650ac7271aa36903c889ca81b4f1f3baa6d4c4918fd4382`.

Lines 103–105:
```rust
    TurnAborted,

    #[error("shared rollout token budget exhausted")]
```

Lines 125–133:
```rust
    #[error(
        "Codex ran out of room in the model's context window. Start a new thread or clear earlier history before retrying."
    )]
    ContextWindowExceeded,
    #[error("no thread with id: {0}")]
    ThreadNotFound(ThreadId),
    #[error("agent thread limit reached")]
    AgentLimitReached { max_threads: usize },
    #[error("session configured event was not the first event in the stream")]
```

Lines 186–191:
```rust
    #[error("Quota exceeded. Check your plan and billing details.")]
    QuotaExceeded,
    #[error(
        "To use Codex with your ChatGPT plan, upgrade to Plus: https://chatgpt.com/explore/plus."
    )]
    UsageNotIncluded,
```

Lines 747–771:
```rust
        if let Some(rate_limit_reached_type) = self.rate_limit_reached_type {
            match rate_limit_reached_type {
                RateLimitReachedType::WorkspaceOwnerCreditsDepleted => {
                    return write!(
                        f,
                        "Your workspace is out of credits. Add credits to continue."
                    );
                }
                RateLimitReachedType::WorkspaceMemberCreditsDepleted => {
                    return write!(
                        f,
                        "Your workspace is out of credits. Ask your workspace owner to refill in order to continue."
                    );
                }
                RateLimitReachedType::WorkspaceOwnerUsageLimitReached => {
                    return write!(
                        f,
                        "You hit your spend cap set in your workspace. Increase your spend cap to continue."
                    );
                }
                RateLimitReachedType::WorkspaceMemberUsageLimitReached => {
                    return write!(
                        f,
                        "You hit your spend cap set by the owner of your workspace. Ask an owner to increase your spend cap to continue."
                    );
```

Lines 820–828:
```rust
            Some(PlanType::Known(
                KnownPlan::Enterprise | KnownPlan::Edu | KnownPlan::EduPlus | KnownPlan::EduPro,
            )) => format!(
                "You’ve hit your usage limit.{}",
                retry_suffix(self.resets_at.as_ref())
            ),
            Some(PlanType::Unknown(_)) | None => format!(
                "You’ve hit your usage limit.{}",
                retry_suffix(self.resets_at.as_ref())
```
