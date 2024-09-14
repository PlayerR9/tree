package Tree

import (
	"errors"

	ud "github.com/PlayerR9/MyGoLib/Units/Debugging"
	uc "github.com/PlayerR9/MyGoLib/Units/common"
	us "github.com/PlayerR9/MyGoLib/Units/slice"
)

// NewTree creates a new tree with the given root.
//
// Parameters:
//   - root: The root of the tree.
//
// Returns:
//   - *Debugging.History[*Tree]: A pointer to the history of the tree.
func NewTreeWithHistory(root *TreeNode[T]) *ud.History[*Tree] {
	tree := NewTree(root)

	h := ud.NewHistory(tree)

	return h
}

// SetChildrenCmd is a command that sets the children of a node.
type SetChildrenCmd struct {
	// prev_tree is a copy of the tree before setting the children.
	prev_tree *Tree

	// children is a slice of pointers to the children to set.
	children []*Tree
}

// Execute implements the Debugging.Commander interface.
func (cmd *SetChildrenCmd) Execute(data *Tree) error {
	cmd.prev_tree = data.Copy().(*Tree)

	err := data.SetChildren(cmd.children)
	if err != nil {
		return err
	}

	return nil
}

// Undo implements the Debugging.Commander interface.
//
// Never errors.
func (cmd *SetChildrenCmd) Undo(data *Tree) error {
	data.root = cmd.prev_tree.root
	data.leaves = cmd.prev_tree.leaves
	data.size = cmd.prev_tree.size

	return nil
}

// Copy implements the Debugging.Commander interface.
func (cmd *SetChildrenCmd) Copy() uc.Copier {
	tree_copy := cmd.prev_tree.Copy().(*Tree)

	children_copy := uc.SliceCopy(cmd.children)

	c_copy := &SetChildrenCmd{
		children:  children_copy,
		prev_tree: tree_copy,
	}

	return c_copy
}

// NewSetChildrenCmd creates a new SetChildrenCmd.
//
// Parameters:
//   - children: The children to set.
//
// Returns:
//   - *SetChildrenCmd: A pointer to the new SetChildrenCmd.
func NewSetChildrenCmd(children []*Tree) *SetChildrenCmd {
	children = us.SliceFilter(children, FilterNonNilTree)
	if len(children) == 0 {
		return nil
	}

	cmd := &SetChildrenCmd{
		children: children,
	}

	return cmd
}

// CleanupCmd is a command that cleans up the tree.
type CleanupCmd struct {
	// root is a pointer to the root of the tree.
	root *TreeNode[T]
}

// Execute implements the Debugging.Commander interface.
//
// Never errors.
func (cmd *CleanupCmd) Execute(data *Tree) error {
	cmd.root = data.root.Copy().(*TreeNode[T])

	data.Cleanup()

	return nil
}

// Undo implements the Debugging.Commander interface.
//
// Never errors.
func (cmd *CleanupCmd) Undo(data *Tree) error {
	data.root = cmd.root

	return nil
}

// Copy implements the Debugging.Commander interface.
func (cmd *CleanupCmd) Copy() uc.Copier {
	root_copy := cmd.root.Copy().(*TreeNode[T])

	cmd_copy := &CleanupCmd{
		root: root_copy,
	}

	return cmd_copy
}

// NewCleanupCmd creates a new CleanupCmd.
//
// Returns:
//   - *CleanupCmd: A pointer to the new CleanupCmd.
func NewCleanupCmd() *CleanupCmd {
	cmd := &CleanupCmd{}
	return cmd
}

// RegenerateLeavesCmd is a command that regenerates the leaves of the tree.
type RegenerateLeavesCmd struct {
	// leaves is a slice of pointers to the leaves of the tree.
	leaves []*TreeNode[T]

	// size is the size of the tree before regenerating the leaves.
	size int
}

// Execute implements the Debugging.Commander interface.
func (cmd *RegenerateLeavesCmd) Execute(data *Tree) error {
	cmd.leaves = data.leaves
	cmd.size = data.size

	data.RegenerateLeaves()

	return nil
}

// Undo implements the Debugging.Commander interface.
func (cmd *RegenerateLeavesCmd) Undo(data *Tree) error {
	data.leaves = cmd.leaves
	data.size = cmd.size

	return nil
}

// Copy implements the Debugging.Commander interface.
func (cmd *RegenerateLeavesCmd) Copy() uc.Copier {
	leaves := make([]*TreeNode[T], len(cmd.leaves))
	copy(leaves, cmd.leaves)

	cmd_copy := &RegenerateLeavesCmd{
		leaves: leaves,
		size:   cmd.size,
	}

	return cmd_copy
}

// NewRegenerateLeavesCmd creates a new RegenerateLeavesCmd.
//
// Returns:
//   - *RegenerateLeavesCmd: A pointer to the new RegenerateLeavesCmd.
func NewRegenerateLeavesCmd() *RegenerateLeavesCmd {
	cmd := &RegenerateLeavesCmd{}
	return cmd
}

// UpdateLeavesCmd is a command that updates the leaves of the tree.
type UpdateLeavesCmd struct {
	// leaves is a slice of pointers to the leaves of the tree.
	leaves []*TreeNode[T]

	// size is the size of the tree before updating the leaves.
	size int
}

// Execute implements the Debugging.Commander interface.
func (cmd *UpdateLeavesCmd) Execute(data *Tree) error {
	cmd.leaves = data.leaves
	cmd.size = data.size

	data.UpdateLeaves()

	return nil
}

// Undo implements the Debugging.Commander interface.
func (cmd *UpdateLeavesCmd) Undo(data *Tree) error {
	data.leaves = cmd.leaves
	data.size = cmd.size

	return nil
}

// Copy implements the Debugging.Commander interface.
func (cmd *UpdateLeavesCmd) Copy() uc.Copier {
	leaves := make([]*TreeNode[T], len(cmd.leaves))
	copy(leaves, cmd.leaves)

	cmd_copy := &UpdateLeavesCmd{
		leaves: leaves,
		size:   cmd.size,
	}

	return cmd_copy
}

// NewUpdateLeavesCmd creates a new UpdateLeavesCmd.
//
// Returns:
//   - *UpdateLeavesCmd: A pointer to the new UpdateLeavesCmd.
func NewUpdateLeavesCmd() *UpdateLeavesCmd {
	cmd := &UpdateLeavesCmd{}
	return cmd
}

// PruneBranchesCmd is a command that prunes the branches of the tree.
type PruneBranchesCmd struct {
	// tree is a pointer to the tree before pruning the branches.
	tree *Tree

	// filter is the filter to apply to prune the branches.
	filter us.PredicateFilter[*TreeNode[T]]

	// ok is true if the whole tree can be deleted, false otherwise.
	ok bool
}

// Execute implements the Debugging.Commander interface.
//
// Never errors.
func (cmd *PruneBranchesCmd) Execute(data *Tree) error {
	cmd.tree = data.Copy().(*Tree)

	cmd.ok = data.PruneBranches(cmd.filter)

	return nil
}

// Undo implements the Debugging.Commander interface.
//
// Never errors.
func (cmd *PruneBranchesCmd) Undo(data *Tree) error {
	data.root = cmd.tree.root
	data.leaves = cmd.tree.leaves
	data.size = cmd.tree.size

	return nil
}

// Copy implements the Debugging.Commander interface.
func (cmd *PruneBranchesCmd) Copy() uc.Copier {
	tree := cmd.tree.Copy().(*Tree)

	cmd_copy := &PruneBranchesCmd{
		tree:   tree,
		filter: cmd.filter,
		ok:     cmd.ok,
	}

	return cmd_copy
}

// NewPruneBranchesCmd creates a new PruneBranchesCmd.
//
// Parameters:
//   - filter: The filter to apply to prune the branches.
//
// Returns:
//   - *PruneBranchesCmd: A pointer to the new PruneBranchesCmd.
func NewPruneBranchesCmd(filter us.PredicateFilter[*TreeNode[T]]) *PruneBranchesCmd {
	if filter == nil {
		return nil
	}

	cmd := &PruneBranchesCmd{
		filter: filter,
	}

	return cmd
}

// GetOk returns the value of the ok field.
//
// Call this function after executing the command.
//
// Returns:
//   - bool: The value of the ok field.
func (cmd *PruneBranchesCmd) GetOk() bool {
	return cmd.ok
}

// SkipFuncCmd is a command that skips the nodes of the tree that
// satisfy the given filter.
type SkipFuncCmd struct {
	// tree is a pointer to the tree before skipping the nodes.
	tree *Tree

	// filter is the filter to apply to skip the nodes.
	filter us.PredicateFilter[*TreeNode[T]]

	// trees is a slice of pointers to the trees obtained after
	// skipping the nodes.
	trees []*Tree
}

// Execute implements the Debugging.Commander interface.
//
// Never errors.
func (cmd *SkipFuncCmd) Execute(data *Tree) error {
	cmd.tree = data.Copy().(*Tree)

	cmd.trees = data.SkipFilter(cmd.filter)

	return nil
}

// Undo implements the Debugging.Commander interface.
//
// Never errors.
func (cmd *SkipFuncCmd) Undo(data *Tree) error {
	data.root = cmd.tree.root
	data.leaves = cmd.tree.leaves
	data.size = cmd.tree.size

	return nil
}

// Copy implements the Debugging.Commander interface.
func (cmd *SkipFuncCmd) Copy() uc.Copier {
	tree := cmd.tree.Copy().(*Tree)

	trees := uc.SliceCopy(cmd.trees)

	cmd_copy := &SkipFuncCmd{
		tree:   tree,
		filter: cmd.filter,
		trees:  trees,
	}

	return cmd_copy
}

// NewSkipFuncCmd creates a new SkipFuncCmd.
//
// Parameters:
//   - filter: The filter to apply to skip the nodes.
//
// Returns:
//   - *SkipFuncCmd: A pointer to the new SkipFuncCmd.
func NewSkipFuncCmd(filter us.PredicateFilter[*TreeNode[T]]) *SkipFuncCmd {
	if filter == nil {
		return nil
	}

	cmd := &SkipFuncCmd{
		filter: filter,
	}

	return cmd
}

// GetTrees returns the value of the trees field.
//
// Call this function after executing the command.
//
// Returns:
//   - []*Tree: A slice of pointers to the trees obtained after
//     skipping the nodes.
func (cmd *SkipFuncCmd) GetTrees() []*Tree {
	return cmd.trees
}

// ProcessLeavesCmd is a command that processes the leaves of the tree.
type ProcessLeavesCmd struct {
	// leaves is a slice of pointers to the leaves of the tree.
	leaves []*TreeNode[T]

	// f is the function to apply to the leaves.
	f uc.EvalManyFunc[*TreeNode[T], *TreeNode[T]]
}

// Execute implements the Debugging.Commander interface.
func (cmd *ProcessLeavesCmd) Execute(data *Tree) error {
	cmd.leaves = data.leaves

	err := data.ProcessLeaves(cmd.f)
	if err != nil {
		return err
	}

	return nil
}

// Undo implements the Debugging.Commander interface.
func (cmd *ProcessLeavesCmd) Undo(data *Tree) error {
	data.leaves = cmd.leaves

	return nil
}

// Copy implements the Debugging.Commander interface.
func (cmd *ProcessLeavesCmd) Copy() uc.Copier {
	leaves := make([]*TreeNode[T], len(cmd.leaves))
	copy(leaves, cmd.leaves)

	cmd_copy := &ProcessLeavesCmd{
		leaves: leaves,
		f:      cmd.f,
	}

	return cmd_copy
}

// NewProcessLeavesCmd creates a new ProcessLeavesCmd.
//
// Parameters:
//   - f: The function to apply to the leaves.
//
// Returns:
//   - *ProcessLeavesCmd: A pointer to the new ProcessLeavesCmd.
func NewProcessLeavesCmd(f uc.EvalManyFunc[*TreeNode[T], *TreeNode[T]]) *ProcessLeavesCmd {
	if f == nil {
		return nil
	}

	cmd := &ProcessLeavesCmd{
		f: f,
	}

	return cmd
}

// DeleteBranchContainingCmd is a command that deletes the branch containing
// the given node.
type DeleteBranchContainingCmd struct {
	// tree is a pointer to the tree before deleting the branch.
	tree *Tree

	// tn is a pointer to the node to delete.
	tn *TreeNode[T]
}

// Execute implements the Debugging.Commander interface.
func (cmd *DeleteBranchContainingCmd) Execute(data *Tree) error {
	cmd.tree = data.Copy().(*Tree)

	err := data.DeleteBranchContaining(cmd.tn)
	if err != nil {
		return err
	}

	return nil
}

// Undo implements the Debugging.Commander interface.
func (cmd *DeleteBranchContainingCmd) Undo(data *Tree) error {
	data.root = cmd.tree.root
	data.leaves = cmd.tree.leaves
	data.size = cmd.tree.size

	return nil
}

// Copy implements the Debugging.Commander interface.
func (cmd *DeleteBranchContainingCmd) Copy() uc.Copier {
	tree := cmd.tree.Copy().(*Tree)
	tn := cmd.tn.Copy().(*TreeNode[T])

	cmd_copy := &DeleteBranchContainingCmd{
		tree: tree,
		tn:   tn,
	}

	return cmd_copy
}

// NewDeleteBranchContainingCmd creates a new DeleteBranchContainingCmd.
//
// Parameters:
//   - tn: The node to delete.
//
// Returns:
//   - *DeleteBranchContainingCmd: A pointer to the new DeleteBranchContainingCmd.
func NewDeleteBranchContainingCmd(tn *TreeNode[T]) *DeleteBranchContainingCmd {
	if tn == nil {
		return nil
	}

	cmd := &DeleteBranchContainingCmd{
		tn: tn,
	}

	return cmd
}

// PruneTreeCmd is a command that prunes the tree using the given filter.
type PruneTreeCmd struct {
	// tree is a pointer to the tree before pruning.
	tree *Tree

	// filter is the filter to use to prune the tree.
	filter us.PredicateFilter[*TreeNode[T]]

	// ok is true if no nodes were pruned, false otherwise.
	ok bool
}

// Execute implements the Debugging.Commander interface.
func (cmd *PruneTreeCmd) Execute(data *Tree) error {
	cmd.tree = data.Copy().(*Tree)

	ok, err := data.Prune(cmd.filter)
	if err != nil {
		return err
	}

	cmd.ok = ok

	return nil
}

// Undo implements the Debugging.Commander interface.
func (cmd *PruneTreeCmd) Undo(data *Tree) error {
	data.root = cmd.tree.root
	data.leaves = cmd.tree.leaves
	data.size = cmd.tree.size

	return nil
}

// Copy implements the Debugging.Commander interface.
func (cmd *PruneTreeCmd) Copy() uc.Copier {
	tree := cmd.tree.Copy().(*Tree)

	cmd_copy := &PruneTreeCmd{
		tree:   tree,
		filter: cmd.filter,
		ok:     cmd.ok,
	}

	return cmd_copy
}

// NewPruneTreeCmd creates a new PruneTreeCmd.
//
// Parameters:
//   - filter: The filter to use to prune the tree.
//
// Returns:
//   - *PruneTreeCmd: A pointer to the new PruneTreeCmd.
func NewPruneTreeCmd(filter us.PredicateFilter[*TreeNode[T]]) *PruneTreeCmd {
	if filter == nil {
		return nil
	}

	cmd := &PruneTreeCmd{
		filter: filter,
	}

	return cmd
}

// GetOk returns the value of the ok field.
//
// Call this function after executing the command.
//
// Returns:
//   - bool: The value of the ok field.
func (cmd *PruneTreeCmd) GetOk() bool {
	return cmd.ok
}

// ExtractBranchCmd is a command that extracts the branch containing the
// given node.
type ExtractBranchCmd struct {
	// leaf is a pointer to the leaf to extract the branch from.
	leaf *TreeNode[T]

	// branch is a pointer to the branch extracted.
	branch *Branch
}

// Execute implements the Debugging.Commander interface.
func (cmd *ExtractBranchCmd) Execute(data *Tree) error {
	branch, err := data.ExtractBranch(cmd.leaf, true)
	if err != nil {
		return err
	}

	cmd.branch = branch

	return nil
}

// Undo implements the Debugging.Commander interface.
func (cmd *ExtractBranchCmd) Undo(data *Tree) error {
	return nil
}

// Copy implements the Debugging.Commander interface.
func (cmd *ExtractBranchCmd) Copy() uc.Copier {
	leaf_copy := cmd.leaf.Copy().(*TreeNode[T])
	branch_copy := cmd.branch.Copy().(*Branch)

	cmd_copy := &ExtractBranchCmd{
		leaf:   leaf_copy,
		branch: branch_copy,
	}

	return cmd_copy
}

// NewExtractBranchCmd creates a new ExtractBranchCmd.
//
// Parameters:
//   - leaf: The leaf to extract the branch from.
//
// Returns:
//   - *ExtractBranchCmd: A pointer to the new ExtractBranchCmd.
func NewExtractBranchCmd(leaf *TreeNode[T]) *ExtractBranchCmd {
	cmd := &ExtractBranchCmd{
		leaf: leaf,
	}
	return cmd
}

// GetBranch returns the value of the branch field.
//
// Call this function after executing the command.
//
// Returns:
//   - *Branch[T]: A pointer to the branch extracted.
func (cmd *ExtractBranchCmd) GetBranch() *Branch {
	branch := cmd.branch
	return branch
}

// InsertBranchCmd is a command that inserts a branch into the tree.
type InsertBranchCmd struct {
	// branch is a pointer to the branch to insert.
	branch *Branch

	// hasError is true if an error occurred during execution, false otherwise.
	hasError bool
}

// Execute implements the Debugging.Commander interface.
func (cmd *InsertBranchCmd) Execute(data *Tree) error {
	ok, err := data.InsertBranch(cmd.branch)
	if err != nil {
		return err
	}

	if !ok {
		cmd.hasError = true
	}

	return nil
}

// Undo implements the Debugging.Commander interface.
func (cmd *InsertBranchCmd) Undo(data *Tree) error {
	err := data.DeleteBranchContaining(cmd.branch.from_node)
	if err != nil {
		if cmd.hasError {
			return nil
		}

		return err
	} else if cmd.hasError {
		return errors.New("error occurred during execution")
	}

	return nil
}

// Copy implements the Debugging.Commander interface.
func (cmd *InsertBranchCmd) Copy() uc.Copier {
	branch_copy := cmd.branch.Copy().(*Branch)

	cmd_copy := &InsertBranchCmd{
		branch:   branch_copy,
		hasError: cmd.hasError,
	}

	return cmd_copy
}

// NewInsertBranchCmd creates a new InsertBranchCmd.
//
// Parameters:
//   - branch: The branch to insert.
//
// Returns:
//   - *InsertBranchCmd: A pointer to the new InsertBranchCmd.
func NewInsertBranchCmd(branch *Branch) *InsertBranchCmd {
	if branch == nil {
		return nil
	}

	cmd := &InsertBranchCmd{
		branch: branch,
	}

	return cmd
}
