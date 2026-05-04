import { Component } from '@angular/core';
import { GraphVisualizationComponent } from './components/graph-visualization/graph-visualization.component';

@Component({
  selector: 'app-root',
  imports: [GraphVisualizationComponent],
  templateUrl: './app.component.html',
  styleUrl: './app.component.css',
})
export class AppComponent {
  title = 'dsr-sim';
}
