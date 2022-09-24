import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { HomeComponent } from './home/home.component';
import { ProfileComponent } from './profile/profile.component';
import { SignInComponent } from './sign-in/sign-in.component';
import { SignUpComponent } from './sign-up/sign-up.component';
import { AuthGuardService } from './auth-guard.service';
import { MqttClientComponent } from './mqtt-client/mqtt-client.component';
import { ParseSignalsComponent } from './parse-signals/parse-signals.component';

const routes: Routes = [

  { path: '', redirectTo: 'signIn', pathMatch: 'full' },
  { path: 'signIn',component: SignInComponent},
  { path: 'signUp',component: SignUpComponent},
  { path: 'profile',component: ProfileComponent, canActivate: [AuthGuardService]},
  { path: 'home', component: HomeComponent, canActivate: [AuthGuardService] },
  { path: 'mqtt-client', component: MqttClientComponent,canActivate: [AuthGuardService] },
  { path: 'parse-signals', component: ParseSignalsComponent,canActivate: [AuthGuardService] },
  { path: '**', redirectTo: 'signIn'},
 

];

@NgModule({
  imports: [RouterModule.forRoot(routes, {onSameUrlNavigation: 'reload'})],
  exports: [RouterModule]
})
export class AppRoutingModule { }
