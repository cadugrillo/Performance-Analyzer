import { Component, OnInit } from '@angular/core';
import { TodoService, Todo } from '../todo.service';
import { CognitoService } from '../cognito.service';


@Component({
  selector: 'app-todo',
  templateUrl: './todo.component.html',
  styleUrls: ['./todo.component.css']
})
export class TodoComponent implements OnInit {

  activeTodos: Todo[] = [];
  completedTodos: Todo[] = [];
  todoMessage!: string;
  userId!: string;
  displayedColumns: string[] = ['description', 'actions']

  constructor(private todoService: TodoService,
              private cognitoService: CognitoService) { }

  ngOnInit() {
    this.cognitoService.getUserId().then((userId: string) => {
      this.userId = userId;
      this.getAll();
    });
  }

  getAll() {
    this.todoService.getTodoList(this.userId).subscribe((data) => {
      this.activeTodos = (data as Todo[]).filter((a) => !a.complete);
      this.completedTodos = (data as Todo[]).filter((a) => a.complete);
    });
  }

  addTodo() {
    var newTodo : Todo = {
      message: this.todoMessage,
      id: '',
      complete: false
    };

    this.todoService.addTodo(this.userId, newTodo).subscribe(() => {
      this.getAll();
      this.todoMessage = '';
    });
  }

  completeTodo(todo: Todo) {
    this.todoService.completeTodo(this.userId, todo).subscribe(() => {
      this.getAll();
    });
  }

  deleteTodo(todo: Todo) {
    this.todoService.deleteTodo(this.userId, todo).subscribe(() => {
      this.getAll();
    })
  }
}
