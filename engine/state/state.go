// これらを持つ
// // the task
// task: Box<dyn Task>,
// // predefined variables
// variables: HashMap<String, String>,
// // model memories, goals and other storages
// storages: HashMap<String, Storage>,
// // available actions and execution history
// namespaces: Vec<Namespace>,
// // list of executed actions
// history: History,
// // optional rag engine
// rag: Option<mini_rag::VectorStore>,
// // set to true when task is complete
// complete: bool,
// // events channel
// events_tx: super::events::Sender,
// // runtime metrics
// pub metrics: Metrics,
// // model support stool
// pub native_tools_support: bool,
package state

type State struct {
}
