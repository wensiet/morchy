import argparse
import asyncio
import logging

from src.settings import AgentSettings, ControlPlaneSettings
from src.scenarios import (
    BasicWorkloadAcquisition,
    MultipleWorkloads,
    Rescheduling,
    ConcurrentWorkloads,
)


logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s | %(message)s",
)
logging.getLogger("httpx").disabled = True
logging.getLogger("httpcore").disabled = True
logging.getLogger("paramiko").disabled = True
logging.getLogger("paramiko.transport").disabled = True

logger = logging.getLogger(__name__)


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="Run workload orchestration scenarios")
    parser.add_argument(
        "-s",
        "--scenario",
        type=str,
        default="ALL",
        help="Scenario name to run (default: ALL)",
    )
    parser.add_argument(
        "-n",
        "--num-trials",
        type=int,
        default=1,
        help="Number of trials to run (default: 1)",
    )
    return parser.parse_args()


async def main() -> None:
    args = parse_args()

    cp_settings = ControlPlaneSettings()
    agent_settings = AgentSettings()

    scenarios = {
        "BasicWorkloadAcquisition": BasicWorkloadAcquisition(
            cp_settings, agent_settings
        ),
        "MultipleWorkloads": MultipleWorkloads(cp_settings, agent_settings),
        "Rescheduling": Rescheduling(cp_settings, agent_settings),
        "ConcurrentWorkloads": ConcurrentWorkloads(cp_settings, agent_settings),
    }

    if args.scenario == "ALL":
        scenarios_to_run = scenarios.values()
    elif args.scenario in scenarios:
        scenarios_to_run = [scenarios[args.scenario]]
    else:
        logger.error(f"No scenario found matching '{args.scenario}'")
        logger.info(f"Available scenarios: {', '.join(scenarios.keys())}")
        return

    for scenario in scenarios_to_run:
        scenario_name = scenario.__class__.__name__
        logger.info(f"Starting {scenario_name} scenario ({args.num_trials} trial(s))")

        for trial in range(args.num_trials):
            logger.info(f"Trial {trial + 1}/{args.num_trials} for {scenario_name}")
            await scenario.run()

        logger.info(f"Completed {scenario_name} scenario")


if __name__ == "__main__":
    asyncio.run(main())
