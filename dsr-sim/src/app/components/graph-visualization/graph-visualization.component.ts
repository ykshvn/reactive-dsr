import {
  Component,
  ViewChild,
  ElementRef,
  OnDestroy,
  AfterViewInit,
} from '@angular/core';
import cytoscape from 'cytoscape';
import {
  GraphResponse,
  GraphService,
} from '../../services/graph.service.ts.service';
import { FormsModule } from '@angular/forms';
import { MatSliderDragEvent, MatSliderModule } from '@angular/material/slider';
import { MatButtonModule } from '@angular/material/button';
import { MatInputModule } from '@angular/material/input';

@Component({
  selector: 'app-graph-visualization',
  templateUrl: './graph-visualization.component.html',
  imports: [FormsModule, MatSliderModule, MatButtonModule, MatInputModule],
  styleUrls: ['./graph-visualization.component.css'],
})
export class GraphVisualizationComponent implements AfterViewInit, OnDestroy {
  @ViewChild('cy') cyContainer!: ElementRef;

  private cy: any;
  graphData: GraphResponse | null = null;
  private ws!: WebSocket;

  thumbLabel = true;

  nodesCount = 20;
  sourceId = 0;
  destId = 10;

  events: any[] = [];

  isRunning: boolean = false;

  constructor(private graphService: GraphService) {}

  ngAfterViewInit(): void {
    this.initCytoscape();
    this.initWebSocket();
    this.loadGraph();
  }

  private initCytoscape() {
    this.cy = cytoscape({
      container: this.cyContainer.nativeElement,
      style: [
        {
          selector: 'node',
          style: {
            'background-color': '#4a90e2',
            label: 'data(id)',
            width: 38,
            height: 38,
            'font-size': '20px',
            color: '#fff',
            'text-valign': 'center',
            'font-weight': 'bold',
          },
        },
        {
          selector: 'edge',
          style: {
            width: 2,
            'line-color': '#888',
            'curve-style': 'bezier',
          },
        },
      ],
      layout: { name: 'preset' },
      userZoomingEnabled: true,
      userPanningEnabled: true,
    });
  }

  private initWebSocket() {
    this.ws = new WebSocket('ws://localhost:6969/ws/simulation');

    this.ws.onopen = () => console.log('WebSocket connected');

    this.ws.onmessage = (event) => {
      const data = JSON.parse(event.data);
      this.events.unshift(data);

      if (data.type === 'rreq_propagated' && data.payload) {
        const p = data.payload;
        this.highlightNode(p.from, '#ffeb3b', 600);
        if (p.route_so_far && p.route_so_far.length > 1) {
          this.highlightEdge(
            p.route_so_far[p.route_so_far.length - 2],
            p.from,
            '#4ade80',
            800,
          );
        }
      }

      if (data.type === 'rrep_received' && data.payload) {
        const p = data.payload;
        this.highlightNode(p.from, '#4ade80', 800);
        if (p.route && p.route.length > 1) {
          this.highlightEdge(p.from, p.route[1], '#22c55e', 1200);
        }
      }
    };

    this.ws.onerror = (e) => console.error('WebSocket error', e);
  }

  loadGraph(nodesCount: number = 20) {
    this.graphService.generateGraph(nodesCount).subscribe({
      next: (data) => {
        this.graphData = data;
        this.renderGraph(data);
      },
      error: (err) => {
        console.error('Failed to load graph', err);
      },
    });
  }

  startRouteDiscovery() {
    if (!this.graphData) return;

    this.events = [];

    this.isRunning = true;

    this.graphService
      .startRouteDiscovery(this.sourceId, this.destId)
      .subscribe({
        next: () =>
          console.log(
            `Route discovery started: ${this.sourceId} → ${this.destId}`,
          ),
        error: (err) => console.error(err),
      });
  }

  nextStep() {
    this.graphService.nextStep().subscribe({
      next: () => console.log('Step executed'),
      error: (err) => console.error(err),
    });
  }

  private renderGraph(graph: GraphResponse) {
    if (!this.cy) return;

    this.cy.elements().remove();

    const elements: any[] = [];

    graph.nodes.forEach((node) => {
      elements.push({
        data: { id: node.id.toString() },
        position: { x: node.x, y: node.y },
      });
    });

    graph.edges.forEach((edge) => {
      elements.push({
        data: {
          id: `e${edge.from}-${edge.to}`,
          source: edge.from.toString(),
          target: edge.to.toString(),
        },
      });
    });

    this.cy.add(elements);

    this.cy
      .layout({
        name: 'preset',
        animate: true,
        duration: 500,
      })
      .run();
  }

  onSliderRelease(event: MatSliderDragEvent) {
    this.loadGraph(event.value);
  }

  highlightNode(
    nodeId: number,
    color: string = '#ffeb3b',
    duration: number = 800,
  ) {
    const node = this.cy.getElementById(nodeId.toString());
    if (node) {
      const originalColor = node.style('background-color');
      node.style('background-color', color);
      setTimeout(() => {
        node.style('background-color', originalColor);
      }, duration);
    }
  }

  highlightEdge(
    from: number,
    to: number,
    color: string = '#4ade80',
    duration: number = 1000,
  ) {
    const edge = this.cy.getElementById(`e${from}-${to}`);
    if (edge) {
      const originalColor = edge.style('line-color');
      edge.style({ 'line-color': color, width: 4 });
      setTimeout(() => {
        edge.style({ 'line-color': originalColor, width: 2.5 });
      }, duration);
    }
  }

  ngOnDestroy() {
    if (this.ws) this.ws.close();
    if (this.cy) this.cy.destroy();
  }
}
