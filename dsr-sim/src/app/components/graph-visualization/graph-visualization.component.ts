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

@Component({
  selector: 'app-graph-visualization',
  templateUrl: './graph-visualization.component.html',
  styleUrls: ['./graph-visualization.component.css'],
})
export class GraphVisualizationComponent implements AfterViewInit, OnDestroy {
  @ViewChild('cy') cyContainer!: ElementRef;

  private cy: any;
  graphData: GraphResponse | null = null;

  constructor(private graphService: GraphService) {}

  ngAfterViewInit(): void {
    this.initCytoscape();
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
            width: 30,
            height: 30,
            'font-size': '12px',
            'text-valign': 'center',
            color: '#fff',
          },
        },
        {
          selector: 'edge',
          style: {
            width: 2,
            'line-color': '#ccc',
            'curve-style': 'bezier',
          },
        },
      ],
      layout: { name: 'cose' },
    });
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

  ngOnDestroy() {
    if (this.cy) this.cy.destroy();
  }
}
