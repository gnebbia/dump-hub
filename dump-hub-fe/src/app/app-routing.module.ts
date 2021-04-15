import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { MainComponent } from './main/main.component';
import { StatusComponent } from './status/status.component';
import { UploadComponent } from './upload/upload.component';

const routes: Routes = [
  { path: '', component: MainComponent },
  { path: 'upload', component: UploadComponent},
  { path: 'status', component: StatusComponent},
  { path: '**', redirectTo: '', pathMatch: 'full' }
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
