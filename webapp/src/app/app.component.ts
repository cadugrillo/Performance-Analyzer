import { Component } from '@angular/core';
import { Router } from '@angular/router';

import { CognitoService } from './services/cognito.service';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent {
  title = 'ui';

  isAuthenticated: boolean;

  constructor(private router: Router,
              private cognitoService: CognitoService) {
    this.isAuthenticated = false;
  }

  public ngOnInit(): void {
    this.cognitoService.isAuthenticated()
    .then((success: boolean) => {
      this.isAuthenticated = success;
      if (this.isAuthenticated) {
        this.router.navigate(['/home']);
      }
    });
  }

  public signOut(): void {
    this.cognitoService.signOut()
    .then(() => {
      this.router.navigate(['/signIn']).then(() => {window.location.reload();});
    });
  }

  openWebPage() {
    window.open('https://github.com/cadugrillo/Performance-Analyzer', '_blank');
  }
}
