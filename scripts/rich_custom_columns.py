from rich.progress import ProgressColumn
from rich.text import Text

class TaskProgressColumn(ProgressColumn):
    """Renders task progress as 'current/total'."""
    def render(self, task) -> Text:
        """Show task progress."""
        return Text(f"{int(task.completed)}/{int(task.total)}", style="progress.percentage") 